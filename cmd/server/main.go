// Сервер.
package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	var (
		config *server.Config
		logger *zap.SugaredLogger
		err    error
	)

	if config, err = server.InitializeConfig(os.Args[1:]); err != nil {
		log.Fatalf("can't parse flags: %v", err)
	}

	if logger, err = ilog.Initialize(config.LogLevel); err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func() {
		if err = logger.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) {
			log.Printf("can't sync logger: %t", err)
		}
	}()

	logger.Infof("config: %+v", config)

	var dbPool *sql.DB
	if config.IsUseSQLDB() {
		var callback func()
		dbPool, callback = initDBConnection(mainCtx, config.DatabaseDriver, config.DatabaseDSN, logger)
		defer callback()
	}

	var (
		collector *service.MetricCollector
		router    *mux.Router
		wg        sync.WaitGroup
	)

	if collector, err = initMetricCollector(mainCtx, config, dbPool, logger); err != nil {
		log.Fatalf("can't initialize collector: %v", err)
	}

	mh := handler.NewMetricHandler(collector, logger)
	ph := handler.NewPingHandler(dbPool, logger)
	router = mux.NewRouter()

	regMiddleware(router, logger, config.HashKey)

	regMetricRoutes(router, mh)
	regPingRoutes(router, ph)

	httpServer := &http.Server{
		Addr:    config.ServerHost,
		Handler: router,
		BaseContext: func(net.Listener) context.Context {
			return mainCtx
		},
	}
	wg.Add(1)
	go func() {
		logger.Debug("http server started")
		if err = httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				logger.Debug("http server stopped")
			} else {
				logger.Errorf("can't start http server: %v", err)
				stop()
			}
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if err = collector.Backup(mainCtx, exitCtx, config.StoreInterval); err != nil {
			logger.Errorf("can't save statistics snapshot: %v", err)
		}
		wg.Done()
	}()

	go func() {
		<-exitCtx.Done()

		logger.Debug("http server shutting down")
		if err = httpServer.Shutdown(context.Background()); err != nil {
			logger.Errorf("can't shutdown http server: %v", err)
		}
	}()

	wg.Wait()

	logger.Debug("program exited")
}

func regMiddleware(router *mux.Router, logger *zap.SugaredLogger, hashKey string) {
	router.Use(
		middleware.WithBufferedWriter(hashKey),
		middleware.WithContentType,
		middleware.CompressHandler(logger),
		middleware.CheckSumRequestHandler(logger, hashKey),
		middleware.WithLogging(logger),
		middleware.FlushBufferedWriter,
	)
}

func regMetricRoutes(router *mux.Router, mh *handler.MetricHandler) {
	router.HandleFunc("/update/{type}/{name}/{value}", mh.CollectHandler).Methods(http.MethodPost)
	router.HandleFunc("/update/", mh.UpdateJSONHandler).Methods(http.MethodPost)
	router.HandleFunc("/value/{type}/{name}", mh.GetValueHandler).Methods(http.MethodGet)
	router.HandleFunc("/value/", mh.GetJSONValueHandler).Methods(http.MethodPost)
	router.HandleFunc("/", mh.GetAllHandler).Methods(http.MethodGet)
	router.HandleFunc("/updates/", mh.UpdatesJSONHandler).Methods(http.MethodPost)
}

func regPingRoutes(router *mux.Router, ph *handler.PingHandler) {
	router.HandleFunc("/ping", ph.Ping).Methods(http.MethodGet)
}

func initMetricCollector(ctx context.Context, config *server.Config, dbPool *sql.DB, logger *zap.SugaredLogger) (*service.MetricCollector, error) {
	var (
		err       error
		ms        service.MetricRepository
		sn        repository.MetricSnapshot
		collector *service.MetricCollector
	)

	if config.Restore {
		sn = snapshot.NewFileMetricSnapshot(config.FileStoragePath, logger)
	}

	if ms, err = initMetricRepository(sn, dbPool, config.IsUseSQLDB(), logger); err != nil {
		return nil, err
	}

	collector = service.NewMetricCollector(ms, logger)

	if config.Restore {
		if err = collector.Restore(ctx); err != nil {
			return nil, err
		}
		logger.Debug("metric collector restored")
	}

	return collector, nil
}

func initMetricRepository(sn repository.MetricSnapshot, dbPool *sql.DB, useSQL bool, logger *zap.SugaredLogger) (service.MetricRepository, error) {
	if useSQL {
		return repository.NewDBMetricRepository(dbPool, sn, logger)
	}

	return repository.NewMemRepository(sn, logger)
}

func initDBConnection(ctx context.Context, driver, dsn string, logger *zap.SugaredLogger) (*sql.DB, func()) {
	var (
		dbPool *sql.DB
		err    error
	)
	if dbPool, err = db.InitializeDB(ctx, driver, dsn, logger); err != nil {
		logger.Fatalf("can't initialize master db: %v", zap.Error(err))
	}

	if err = db.CreateStructure(ctx, dbPool); err != nil {
		logger.Fatalf("can't create structure: %v", zap.Error(err))
	}

	logger.Debug("initDBConnection finished")

	return dbPool, func() {
		if err = db.CloseDB(dbPool); err != nil {
			logger.Errorf("can't close master db: %v", err)
			return
		}
		logger.Debug("close master db successfully")
	}
}
