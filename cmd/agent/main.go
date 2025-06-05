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

	cl := collector.NewRuntimeMetricCollector()
	clPoller := collector.NewIntervalPoller(cl, time.Duration(config.PollInterval)*time.Second, logger)

	t := transport.NewHTTPClient(config.ServerProtocol+"://"+config.ServerHost, logger)
	sn := sender.NewMetricSender(t, config.BatchEnabled)

	statService := service.NewStatService(cl, sn, time.Duration(config.ReportInterval)*time.Second, logger)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		clPoller.PollStat(ctx)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		statService.SendStat(ctx)
		wg.Done()
	}()

	wg.Wait()

	logger.Debug("program exited")
}
