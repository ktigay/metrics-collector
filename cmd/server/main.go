// Сервер.
package main

import (
	"context"
	"errors"
	"log"
	"net"
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
	"github.com/ktigay/metrics-collector/internal/server/handler"
	"github.com/ktigay/metrics-collector/internal/server/middleware"
	"github.com/ktigay/metrics-collector/internal/server/repository"
	"github.com/ktigay/metrics-collector/internal/server/service"
	"github.com/ktigay/metrics-collector/internal/server/snapshot"
	"go.uber.org/zap"
)

func main() {
	mainCtx := context.TODO()
	exitCtx, stop := signal.NotifyContext(mainCtx, os.Interrupt, syscall.SIGTERM)
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

	ilog.AppLogger.Infof("config: %+v", config)

	if config.IsUseSQLDB() {
		callback := initDBConnection(mainCtx, config.DatabaseDriver, config.DatabaseDSN)
		defer callback()
	}

	var (
		collector *service.MetricCollector
		router    *mux.Router
		mh        *handler.MetricHandler
		wg        sync.WaitGroup
	)

	if collector, err = initMetricCollector(mainCtx, config); err != nil {
		log.Fatalf("can't initialize collector: %v", err)
	}

	mh = handler.NewMetricHandler(collector)
	router = mux.NewRouter()

	registerMiddleware(router)
	registerRoutes(router, mh)

	httpServer := &http.Server{
		Addr:    config.ServerHost,
		Handler: router,
		BaseContext: func(net.Listener) context.Context {
			return mainCtx
		},
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
		if err = saveSnapshot(mainCtx, exitCtx, collector, config.StoreInterval); err != nil {
			ilog.AppLogger.Errorf("can't save statistics snapshot: %v", err)
		}
		wg.Done()
	}()

	go func() {
		<-exitCtx.Done()

		ilog.AppLogger.Debug("http server shutting down")
		if err = httpServer.Shutdown(context.Background()); err != nil {
			ilog.AppLogger.Errorf("can't shutdown http server: %v", err)
		}
	}()

	wg.Wait()

	ilog.AppLogger.Debug("program exited")
}

func registerMiddleware(router *mux.Router) {
	router.Use(
		middleware.WithContentType,
		middleware.CompressHandler,
		middleware.WithLogging,
	)
}

func registerRoutes(router *mux.Router, mh *handler.MetricHandler) {
	router.HandleFunc("/update/{type}/{name}/{value}", mh.CollectHandler).Methods(http.MethodPost)
	router.HandleFunc("/update/", mh.UpdateJSONHandler).Methods(http.MethodPost)
	router.HandleFunc("/value/{type}/{name}", mh.GetValueHandler).Methods(http.MethodGet)
	router.HandleFunc("/value/", mh.GetJSONValueHandler).Methods(http.MethodPost)
	router.HandleFunc("/ping", handler.PingHandler).Methods(http.MethodGet)
	router.HandleFunc("/", mh.GetAllHandler).Methods(http.MethodGet)
	router.HandleFunc("/updates/", mh.UpdatesJSONHandler).Methods(http.MethodPost)
}

func initMetricCollector(ctx context.Context, config *server.Config) (*service.MetricCollector, error) {
	var (
		err       error
		ms        service.MetricRepository
		sn        repository.MetricSnapshot
		collector *service.MetricCollector
	)

	if config.Restore {
		sn = snapshot.NewFileMetricSnapshot(config.FileStoragePath)
	}

	if ms, err = initMetricRepository(sn, config.IsUseSQLDB()); err != nil {
		return nil, err
	}

	collector = service.NewMetricCollector(ms)

	if config.Restore {
		if err = collector.Restore(ctx); err != nil {
			return nil, err
		}
		ilog.AppLogger.Debug("metric collector restored")
	}

	return collector, nil
}

func initMetricRepository(sn repository.MetricSnapshot, useSQL bool) (service.MetricRepository, error) {
	if useSQL {
		return repository.NewDBMetricRepository(db.MasterDB, sn)
	}

	return repository.NewMemRepository(sn)
}

func saveSnapshot(mainCtx, exitCtx context.Context, c *service.MetricCollector, storeInterval int) error {
	ticker := time.NewTicker(time.Duration(storeInterval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ilog.AppLogger.Debug("saveSnapshot saving metrics started")
			if err := c.Backup(mainCtx); err != nil {
				return err
			}
			ilog.AppLogger.Debug("saveSnapshot saving metrics finished")
		case <-exitCtx.Done():
			ticker.Stop()
			if err := c.Backup(mainCtx); err != nil {
				return err
			}
			ilog.AppLogger.Debug("saveSnapshot shutting down")
			return nil
		}
	}
}

func initDBConnection(ctx context.Context, driver, dsn string) func() {
	var err error
	if err = db.InitializeMasterDB(ctx, driver, dsn); err != nil {
		ilog.AppLogger.Fatalf("can't initialize master db: %v", zap.Error(err))
	}

	if err = db.CreateStructure(ctx); err != nil {
		ilog.AppLogger.Fatalf("can't create structure: %v", zap.Error(err))
	}

	ilog.AppLogger.Debug("initDBConnection finished")

	return func() {
		if err = db.CloseMasterDB(); err != nil {
			ilog.AppLogger.Errorf("can't close master db: %v", err)
			return
		}
		ilog.AppLogger.Debug("close master db successfully")
	}
}
