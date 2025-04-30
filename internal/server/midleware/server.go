package midleware

import (
	"github.com/ktigay/metrics-collector/internal/server"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// WithContentType - устанавливает в ResponseWriter Content-Type.
func WithContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// WithLogging - логирует запрос.
func WithLogging(l *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sugar := l.Sugar()

		rd := &server.ResponseData{
			Status: 0,
			Size:   0,
		}
		lw := server.NewResponseWriter(w, rd)
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
