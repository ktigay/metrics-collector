package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"

	"github.com/ktigay/metrics-collector/internal/compress"
	serverhttp "github.com/ktigay/metrics-collector/internal/http"
	"github.com/ktigay/metrics-collector/internal/log"
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
func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rd := &serverhttp.ResponseData{
			Status: 0,
			Size:   0,
		}
		lw := serverhttp.NewWriter(w, rd)

		log.AppLogger.Infow(
			"request",
			"requestURI", r.RequestURI,
			"method", r.Method,
		)

		next.ServeHTTP(lw, r)

		duration := time.Since(start)
		log.AppLogger.Infow(
			"response",
			"status", rd.Status,
			"size", rd.Size,
			"duration", duration,
		)
	})
}

// CompressHandler обработчик сжатия данных.
func CompressHandler(next http.Handler) http.Handler {
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
						log.AppLogger.Error("middleware.CompressHandler error", zap.Error(err))
					}
				}()
			}
		}

		next.ServeHTTP(w, r)
	})
}
