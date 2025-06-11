// Агент.
package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ktigay/metrics-collector/internal/client"
	"github.com/ktigay/metrics-collector/internal/client/collector"
	"github.com/ktigay/metrics-collector/internal/client/sender"
	"github.com/ktigay/metrics-collector/internal/client/sender/transport"
	"github.com/ktigay/metrics-collector/internal/client/service"
	ilog "github.com/ktigay/metrics-collector/internal/log"
	"go.uber.org/zap"
)

func main() {
	var (
		config *client.Config
		logger *zap.SugaredLogger
		err    error
	)

	if config, err = client.InitializeConfig(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	if logger, err = ilog.Initialize(config.LogLevel); err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func() {
		if err = logger.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) {
			log.Printf("can't sync logger: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	storage := collector.NewStorage(logger)
	cl := collector.NewRuntimeMetricCollector(storage)
	rnPoller := collector.NewIntervalPoller(cl, time.Duration(config.PollInterval)*time.Second, logger)

	gp := collector.NewGopsUtilCollector(storage)
	gpPoller := collector.NewIntervalPoller(gp, time.Duration(config.PollInterval)*time.Second, logger)

	t := transport.NewHTTPClient(config.ServerProtocol+"://"+config.ServerHost, config.HashKey, logger)
	sn := sender.NewMetricSender(t, config.BatchEnabled, config.RateLimit, logger)

	statService := service.NewStatService(storage, sn, time.Duration(config.ReportInterval)*time.Second, logger)

	var wg sync.WaitGroup
	tasks := []func(context.Context){
		rnPoller.PollStat,
		gpPoller.PollStat,
		statService.SendStat,
	}

	for _, task := range tasks {
		wg.Add(1)
		go func() {
			task(ctx)
			wg.Done()
		}()
	}

	wg.Wait()

	logger.Debug("program exited")
}
