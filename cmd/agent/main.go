// Агент.
package main

import (
	"context"
	"errors"
	"log"
	"math"
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
	"github.com/ktigay/metrics-collector/internal/metric"
	"go.uber.org/zap"
)

type Task func(context.Context)

func main() {
	var (
		cfg    *client.Config
		logger *zap.SugaredLogger
		err    error
	)

	if cfg, err = client.InitializeConfig(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	if logger, err = ilog.Initialize(cfg.LogLevel); err != nil {
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
	rnPoller := collector.NewIntervalPoller(cl, time.Duration(cfg.PollInterval)*time.Second, logger)

	gp := collector.NewGopsUtilCollector()
	gpPoller := collector.NewIntervalPoller(gp, time.Duration(cfg.PollInterval)*time.Second, logger)

	t := transport.NewHTTPClient(cfg.ServerProtocol+"://"+cfg.ServerHost, cfg.HashKey, logger)
	sn := sender.NewMetricSender(t, cfg.BatchEnabled, cfg.RateLimit, logger)
	handler := collector.NewMetricsHandler()
	statSender := service.NewStatSenderService(sn, handler, time.Duration(cfg.ReportInterval)*time.Second, logger)

	// размер канала такой, чтобы не блокировать сборку статистики.
	chSize := int64(math.Ceil(float64(cfg.ReportInterval)/float64(cfg.PollInterval))) * 2
	pollChan := make(chan []metric.Metrics, chSize)
	defer close(pollChan)

	tasks := []Task{
		func(ctx context.Context) {
			rnPoller.PollStat(ctx, pollChan)
		},
		func(ctx context.Context) {
			gpPoller.PollStat(ctx, pollChan)
		},
		func(ctx context.Context) {
			statSender.SendStat(ctx, pollChan)
		},
	}
	var wg sync.WaitGroup
	wg.Add(len(tasks))

	for _, task := range tasks {
		go func() {
			task(ctx)
			defer wg.Done()
		}()
	}

	wg.Wait()
	logger.Debug("program exited")
}
