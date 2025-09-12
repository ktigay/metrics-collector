// Агент.
package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/ktigay/metrics-collector/internal/client"
	"github.com/ktigay/metrics-collector/internal/client/collector"
	"github.com/ktigay/metrics-collector/internal/client/sender"
	"github.com/ktigay/metrics-collector/internal/client/sender/transport"
	"github.com/ktigay/metrics-collector/internal/client/service"
	"github.com/ktigay/metrics-collector/internal/crypto"
	ilog "github.com/ktigay/metrics-collector/internal/log"
	"github.com/ktigay/metrics-collector/internal/metric"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

// Task задача для запуска в горутинах.
type Task func(context.Context)

func main() {
	var (
		cfg    *client.Config
		logger *zap.SugaredLogger
		err    error
	)

	if err = buildInfo(); err != nil {
		log.Printf("cannot print build info: %s", err)
	}

	if cfg, err = client.InitializeConfig(os.Args[1:]); err != nil {
		handleExit(1)
		return
	}

	if logger, err = ilog.Initialize(cfg.LogLevel); err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func() {
		if err = logger.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) {
			log.Printf("can't sync logger: %v", err)
		}
	}()

	logger.Infof("cfg: %+v", cfg)

	var cryptoKey *crypto.PublicKey
	if cfg.CryptoKey != "" {
		if cryptoKey, err = initPublicKey(cfg.CryptoKey); err != nil {
			log.Fatalf("can't initialize crypto key: %v", err)
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	cl := collector.NewRuntimeMetricCollector()
	rnPoller := collector.NewIntervalPoller(cl, time.Duration(cfg.PollInterval)*time.Second, logger)

	gp := collector.NewGopsUtilCollector()
	gpPoller := collector.NewIntervalPoller(gp, time.Duration(cfg.PollInterval)*time.Second, logger)

	t := getHTTPTransport(cfg.ServerProtocol+"://"+cfg.ServerHost, cfg.HashKey, cryptoKey, logger)
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

func initPublicKey(cryptoKey string) (*crypto.PublicKey, error) {
	file, err := os.Open(cryptoKey)
	if err != nil {
		return nil, err
	}
	return crypto.NewPublicKey(bufio.NewReader(file))
}

func handleExit(code int) {
	os.Exit(code)
}

func buildInfo() error {
	_, err := fmt.Fprintf(os.Stdout, `Build version: %s
Build date: %s
Build commit: %s
`, buildVersion, buildDate, buildCommit)
	return err
}

func getHTTPTransport(
	url,
	hashKey string,
	encryptKey *crypto.PublicKey,
	logger *zap.SugaredLogger,
) sender.Transport {
	factory := transport.NewRequestFactory(http.MethodPost, url, hashKey, encryptKey)
	return transport.NewHTTPClient(factory, logger)
}
