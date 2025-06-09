// Package middleware миддлвары.
package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
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

func WithBufferedWriter(hashKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rd := serverhttp.ResponseData{
				Status: 0,
				Size:   0,
			}

			sw := serverhttp.NewWriter(w, &rd, hashKey)
			next.ServeHTTP(sw, r)
		})
	}
}

func FlushBufferedWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		if bw, ok := w.(*serverhttp.Writer); ok {
			bw.Flush()
			if err := bw.ResponseData().Err; err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	})
}

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

			b, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(b))

			logger.Infow(
				"request",
				"requestURI", r.RequestURI,
				"method", r.Method,
				"body", string(b),
				"headers", r.Header,
			)

			next.ServeHTTP(w, r)

			sw, ok := w.(*serverhttp.Writer)
			if !ok {
				return
			}

			duration := time.Since(start)
			rd := sw.ResponseData()
			logger.Infow(
				"response",
				"duration", duration,
				"status", rd.Status,
				"size", rd.Size,
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

			var aeAlg compress.Type
			if isAccepted {
				aeAlg = compress.TypeFromString(acceptEncoding)
			}

			if string(aeAlg) == "" {
				next.ServeHTTP(w, r)
				return
			}

			var (
				cw  *compress.HTTPWriter
				err error
			)
			if sw, ok := w.(*serverhttp.Writer); ok {
				sw.WithWriter(func(writer http.ResponseWriter) http.ResponseWriter {
					cw, err = compress.NewHTTPWriter(aeAlg, writer)
					if err == nil {
						return cw
					}
					return nil
				})
			} else {
				cw, err = compress.NewHTTPWriter(aeAlg, w)
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r)

			defer func() {
				if err = cw.Close(); err != nil {
					logger.Error("middleware.CompressHandler error", zap.Error(err))
				}
			}()
		})
	}
}

func CheckSumRequestHandler(logger *zap.SugaredLogger, hashKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				checkSum string
				err      error
				buff     []byte
				chBytes  []byte
			)
			if checkSum = r.Header.Get(serverhttp.HashSHA256Header); checkSum == "" {
				next.ServeHTTP(w, r)
				return
			}

			if buff, err = io.ReadAll(r.Body); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(buff))

			if chBytes, err = hex.DecodeString(checkSum); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			b := sha256.Sum256(append(buff, hashKey...))
			if !bytes.Equal(b[:], chBytes) {
				w.WriteHeader(http.StatusBadRequest)
				logger.Warnf("CheckSumRequestHandler: invalid checksum %s", checkSum)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
