package main

import (
	"errors"
	"github.com/ktigay/metrics-collector/internal/server/snapshot"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/ktigay/metrics-collector/internal/server"
	"github.com/ktigay/metrics-collector/internal/server/collector"
	"github.com/ktigay/metrics-collector/internal/server/logger"
	"github.com/ktigay/metrics-collector/internal/server/middleware"
	"github.com/ktigay/metrics-collector/internal/server/storage"
)

func init() {
	path := "./cache"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}

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

	c, err := initCollector(config)
	if err != nil {
		log.Fatalf("can't initialize collector: %v", err)
	}
	s := server.NewServer(c, l)

	router := mux.NewRouter()

	registerMiddleware(router, l)
	registerRoutes(router, s)

	stop := make(chan bool)
	defer func() {
		stop <- true
	}()
	go func() {
		if err = saveStatisticsSnapshot(stop, config, c); err != nil {
			l.Sugar().Errorf("can't save statistics: %v", err)
		}
	}()

	if err = http.ListenAndServe(config.ServerHost, router); err != nil {
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

func initCollector(config *Config) (*collector.MetricCollector, error) {
	var m []storage.Entity

	if config.Restore {
		s, err := snapshot.FileReadAll[storage.Entity](config.FileStoragePath)
		if err != nil {
			return nil, err
		}
		m = s
	}

	return collector.NewMetricCollector(storage.NewMemStorage(m)), nil
}

func saveStatisticsSnapshot(stop <-chan bool, config *Config, c *collector.MetricCollector) error {
	save := func() error {
		return snapshot.FileWriteAll(config.FileStoragePath, c.GetAll())
	}

	for {
		select {
		default:
			time.Sleep(time.Duration(config.StoreInterval) * time.Second)
			if err := save(); err != nil {
				return err
			}
		case <-stop:
			if err := save(); err != nil {
				return err
			}
			return nil
		}
	}
}
