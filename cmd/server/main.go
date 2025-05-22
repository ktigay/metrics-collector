package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
	ilog "github.com/ktigay/metrics-collector/internal/log"
	"github.com/ktigay/metrics-collector/internal/server"
	"github.com/ktigay/metrics-collector/internal/server/db"
	"github.com/ktigay/metrics-collector/internal/server/middleware"
	"github.com/ktigay/metrics-collector/internal/server/service"
	"github.com/ktigay/metrics-collector/internal/server/snapshot"
	"github.com/ktigay/metrics-collector/internal/server/storage"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config, err := server.InitializeConfig(os.Args[1:])
	if err != nil {
		log.Fatalf("can't parse flags: %v", err)
	}

	if err = ilog.Initialize(config.LogLevel); err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func() {
		if err = ilog.AppLogger.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) {
			log.Printf("can't sync logger: %t", err)
		}
	}()

	if err = db.InitializeMasterDB("pgx", config.DatabaseDSN); err != nil {
		log.Fatalf("can't initialize master db: %v", err)
	}

	var (
		collector *service.MetricCollector
		router    *mux.Router
		srv       *server.Server
		wg        sync.WaitGroup
	)

	collector, err = initCollector(config.Restore, config.FileStoragePath)
	if err != nil {
		log.Fatalf("can't initialize collector: %v", err)
	}

	srv = server.NewServer(collector)
	router = mux.NewRouter()

	registerMiddleware(router)
	registerRoutes(router, srv)

	httpServer := &http.Server{
		Addr:    config.ServerHost,
		Handler: router,
	}
	wg.Add(1)
	go func() {
		ilog.AppLogger.Debug("http server started")
		if err = httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				ilog.AppLogger.Debug("http server stopped")
			} else {
				ilog.AppLogger.Errorf("can't start http server: %v", err)
				stop()
			}
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if err = saveSnapshot(ctx, collector, config.StoreInterval); err != nil {
			ilog.AppLogger.Errorf("can't save statistics snapshot: %v", err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		<-ctx.Done()
		ilog.AppLogger.Debug("http server shutting down")
		if err = httpServer.Shutdown(context.Background()); err != nil {
			ilog.AppLogger.Errorf("can't shutdown http server: %v", err)
		}
		if err = db.MasterDB.Close(); err != nil {
			ilog.AppLogger.Errorf("can't close master db: %v", err)
		}
		wg.Done()
	}()

	wg.Wait()

	ilog.AppLogger.Debug("program exited")
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
	router.HandleFunc("/ping", s.Ping).Methods(http.MethodGet)
	router.HandleFunc("/", s.GetAllHandler).Methods(http.MethodGet)
}

func initCollector(restore bool, restorePath string) (*service.MetricCollector, error) {
	var sn storage.Snapshot

	if restore {
		sn = snapshot.NewFileSnapshot(restorePath)
	}

	st, err := storage.NewMemStorage(sn)
	if err != nil {
		return nil, err
	}

	return service.NewMetricCollector(st), nil
}

func saveSnapshot(ctx context.Context, c *service.MetricCollector, storeInterval int) error {
	ticker := time.NewTicker(time.Duration(storeInterval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ilog.AppLogger.Debug("saveSnapshot saving metrics started")
			if err := c.Backup(); err != nil {
				return err
			}
			ilog.AppLogger.Debug("saveSnapshot saving metrics finished")
		case <-ctx.Done():
			ticker.Stop()
			if err := c.Backup(); err != nil {
				return err
			}
			ilog.AppLogger.Debug("saveSnapshot shutting down")
			return nil
		}
	}
}
