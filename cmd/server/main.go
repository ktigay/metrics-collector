package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	ilog "github.com/ktigay/metrics-collector/internal/log"
	"github.com/ktigay/metrics-collector/internal/server"
	"github.com/ktigay/metrics-collector/internal/server/middleware"
	"github.com/ktigay/metrics-collector/internal/server/snapshot"
	"github.com/ktigay/metrics-collector/internal/server/storage"
)

func main() {
	config, err := parseFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("can't parse flags: %v", err)
	}

	if err = ilog.Initialize(config.LogLevel); err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func() {
		if err := ilog.AppLogger.Sync(); err != nil {
			log.Printf("can't sync logger: %v", err)
		}
	}()

	c, err := initCollector(config.Restore, config.FileStoragePath)
	if err != nil {
		log.Fatalf("can't initialize collector: %v", err)
	}
	s := server.NewServer(c)

	router := mux.NewRouter()

	registerMiddleware(router)
	registerRoutes(router, s)

	var wg sync.WaitGroup
	stop := make(chan bool)
	wg.Add(1)
	go func() {
		if err = saveStatisticsSnapshot(stop, c, config.StoreInterval); err != nil {
			ilog.AppLogger.Errorf("can't save statistics: %v", err)
		}
		wg.Done()
	}()

	if err = http.ListenAndServe(config.ServerHost, router); err != nil {
		ilog.AppLogger.Errorln("can't start http server:", err)
	}

	// сигнал к остановке горутины
	stop <- true
	// ждем остановки
	wg.Wait()
}

func registerMiddleware(router *mux.Router) {
	router.Use(
		middleware.WithContentType,
		middleware.WithLogging,
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

func initCollector(restore bool, restorePath string) (*storage.MetricCollector, error) {
	var sn storage.Snapshot

	if restore {
		sn = snapshot.NewFileSnapshot(restorePath)
	}

	st, err := storage.NewMemStorage(sn)
	if err != nil {
		return nil, err
	}

	return storage.NewMetricCollector(st), nil
}

func saveStatisticsSnapshot(stop <-chan bool, c *storage.MetricCollector, storeInterval int) error {
	ticker := time.NewTicker(time.Duration(storeInterval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := c.Backup(); err != nil {
				return err
			}
		case <-stop:
			if err := c.Backup(); err != nil {
				return err
			}
			return nil
		}
	}
}
