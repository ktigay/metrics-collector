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
	ilog "github.com/ktigay/metrics-collector/internal/log"
)

func main() {
	var (
		config *client.Config
		err    error
	)

	if config, err = client.InitializeConfig(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	if err = ilog.Initialize(config.LogLevel); err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func() {
		if err = ilog.AppLogger.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) {
			log.Printf("can't sync logger: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cl := collector.NewRuntimeMetricCollector()
	s := client.NewSender(config.ServerProtocol + "://" + config.ServerHost)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		pollStat(ctx, config, cl)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		sendStat(ctx, config, cl, s)
		wg.Done()
	}()

	wg.Wait()

	ilog.AppLogger.Debug("program exited")
}

func pollStat(ctx context.Context, config *client.Config, cl *collector.RuntimeMetricCollector) {
	ticker := time.NewTicker(time.Duration(config.PollInterval) * time.Second)
	defer ticker.Stop()
	cl.PollStat()
	for {
		select {
		case <-ticker.C:
			cl.PollStat()
		case <-ctx.Done():
			ilog.AppLogger.Debug("pollStat stopped")
			return
		}
	}
}

func sendStat(ctx context.Context, config *client.Config, cl *collector.RuntimeMetricCollector, s *client.Sender) {
	ticker := time.NewTicker(time.Duration(config.ReportInterval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			ilog.AppLogger.Debug("sendStat stopped")
			return
		case <-ticker.C:
			ilog.AppLogger.Debug("sendStat started")
			s.SendMetrics(cl.GetStat())
			ilog.AppLogger.Debug("sendStat finished")
		}
	}
}
