package main

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"os"

	"github.com/ktigay/metrics-collector/internal/server"
	"github.com/ktigay/metrics-collector/internal/server/collector"
	"github.com/ktigay/metrics-collector/internal/server/midleware"
	"github.com/ktigay/metrics-collector/internal/server/storage"

	"net/http"
)

func main() {
	l, err := newLogger("info")
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func() {
		if err := l.Sync(); err != nil {
			log.Fatalf("can't sync logger: %v", err)
		}
	}()

	var config *Config
	if config, err = parseFlags(os.Args[1:]); err != nil {
		l.Sugar().Errorln("can't parse flags: ", err)
	}

	c := collector.NewMetricCollector(storage.NewMemStorage())
	s := server.NewServer(c, l)

	router := mux.NewRouter()

	router.Use(
		midleware.WithContentType,
		func(next http.Handler) http.Handler {
			return midleware.WithLogging(l, next)
		},
	)

	router.HandleFunc("/update/{type}/{name}/{value}", s.CollectHandler).Methods(http.MethodPost)
	router.HandleFunc("/value/{type}/{name}", s.GetValueHandler).Methods(http.MethodGet)
	router.HandleFunc("/", s.GetAllHandler).Methods(http.MethodGet)

	if err := http.ListenAndServe(config.ServerHost, router); err != nil {
		l.Sugar().Errorln("can't start http server:", err)
	}
}

func newLogger(level string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	return cfg.Build()
}
