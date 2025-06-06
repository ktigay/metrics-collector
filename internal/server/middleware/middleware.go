// Package middleware миддлвары.
package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/ktigay/metrics-collector/internal/compress"
	serverhttp "github.com/ktigay/metrics-collector/internal/http"
	"go.uber.org/zap"
)

var acceptTypes = []string{"text/html", "application/json", "*/*"}

// WithContentType устанавливает в ResponseWriter Content-Type.
func WithContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// WithLogging логирует запрос.
func WithLogging(logger *zap.SugaredLogger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rd := serverhttp.ResponseData{
				Status: 0,
				Size:   0,
			}
			lw := serverhttp.NewWriter(w, &rd)

			b, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(b))

			logger.Infow(
				"request",
				"requestURI", r.RequestURI,
				"method", r.Method,
				"body", string(b),
			)

			next.ServeHTTP(lw, r)

			duration := time.Since(start)

			logger.Infow(
				"response",
				"status", rd.Status,
				"size", rd.Size,
				"duration", duration,
				"body", string(rd.Body),
			)
		})
	}
}

// CompressHandler обработчик сжатия данных.
func CompressHandler(logger *zap.SugaredLogger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contentEncoding := r.Header.Get("Content-Encoding")
			if ceAlg := compress.TypeFromString(contentEncoding); ceAlg != "" {
				cr, err := compress.ReaderFactory(ceAlg, r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = cr
			}

			acceptEncoding := r.Header.Get("Accept-Encoding")
			accept := r.Header.Get("Accept")
			isAccepted := func() bool {
				for _, acceptType := range acceptTypes {
					if strings.Contains(accept, acceptType) {
						return true
					}
				}
				return false
			}()

			if isAccepted {
				if aeAlg := compress.TypeFromString(acceptEncoding); string(aeAlg) != "" {
					cw, err := compress.NewHTTPWriter(aeAlg, w)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					w = cw
					defer func() {
						if err = cw.Close(); err != nil {
							logger.Error("middleware.CompressHandler error", zap.Error(err))
						}
					}()
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
