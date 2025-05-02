package middleware

import (
	serverhttp "github.com/ktigay/metrics-collector/internal/server/http"
	serverio "github.com/ktigay/metrics-collector/internal/server/io"
	"strings"

	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	gzipCompression = "gzip"
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
func WithLogging(l *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sugar := l.Sugar()

		rd := &serverhttp.ResponseData{
			Status: 0,
			Size:   0,
		}
		lw := serverhttp.NewWriter(w, rd)
		next.ServeHTTP(lw, r)

		duration := time.Since(start)

		sugar.Infow(
			"request",
			"requestURI", r.RequestURI,
			"method", r.Method,
			"duration", duration,
		)

		sugar.Infow(
			"response",
			"status", rd.Status,
			"size", rd.Size,
		)
	})
}

// CompressHandler обработчик сжатия данных.
func CompressHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		contentEncoding := r.Header.Get("Content-Encoding")
		if strings.Contains(contentEncoding, gzipCompression) {
			cr, err := serverio.CompressReaderFactory(gzipCompression, r.Body)
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
		if isAccepted && strings.Contains(acceptEncoding, gzipCompression) {
			cw, cb := serverhttp.CompressWriterFactory(gzipCompression, w)
			w = cw
			defer func() {
				_ = cb()
			}()
		}

		next.ServeHTTP(w, r)
	})
}
