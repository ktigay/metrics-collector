package main

import (
	"github.com/gorilla/mux"
	"github.com/ktigay/metrics-collector/internal/server"
	"github.com/ktigay/metrics-collector/internal/server/collector"
	"github.com/ktigay/metrics-collector/internal/server/logger"
	"github.com/ktigay/metrics-collector/internal/server/middleware"
	"github.com/ktigay/metrics-collector/internal/server/storage"
	"go.uber.org/zap"
	"log"
	"os"

	"net/http"
)

func main() {
	config, err := parseFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("can't parse flags: %v", err)
	}

	l, err := logger.Initialize(config.LogLevel)
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func() {
		if err := l.Sync(); err != nil {
			log.Fatalf("can't sync logger: %v", err)
		}
	}()

	c := collector.NewMetricCollector(storage.NewMemStorage())
	s := server.NewServer(c, l)

	router := mux.NewRouter()

	registerMiddleware(router, l)
	registerRoutes(router, s)

	if err := http.ListenAndServe(config.ServerHost, router); err != nil {
		l.Sugar().Errorln("can't start http server:", err)
	}
}

func registerMiddleware(router *mux.Router, l *zap.Logger) {
	router.Use(
		middleware.WithContentType,
		func(next http.Handler) http.Handler {
			return middleware.WithLogging(l, next)
		},
		middleware.CompressHandler,
	)
}

func registerRoutes(router *mux.Router, s *server.Server) {
	router.HandleFunc("/update/{type}/{name}/{value}", s.CollectHandler).Methods(http.MethodPost)
	router.HandleFunc("/update/", s.UpdateJSONHandler).Methods(http.MethodPost)
	router.HandleFunc("/value/{type}/{name}", s.GetValueHandler).Methods(http.MethodGet)
	router.HandleFunc("/value/", s.GetJSONValueHandler).Methods(http.MethodPost)
	router.HandleFunc("/", s.GetAllHandler).Methods(http.MethodGet)
}
