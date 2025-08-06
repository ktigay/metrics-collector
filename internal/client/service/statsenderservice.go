package service

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/retry"
)

// MetricsHandler интерфейс обработчика метрик.
//
//go:generate mockgen -destination=./mocks/mock_handler.go -package=mocks github.com/ktigay/metrics-collector/internal/client/service MetricsHandler
type MetricsHandler interface {
	Processing(metrics [][]metric.Metrics) []metric.Metrics
}

// StatSender отправка метрик.
//
//go:generate mockgen -destination=./mocks/mock_sender.go -package=mocks github.com/ktigay/metrics-collector/internal/client/service StatSender
type StatSender interface {
	SendMetrics([]metric.Metrics, chan<- error)
}

// StatSenderService провайдер статистики.
type StatSenderService struct {
	sender   StatSender
	handler  MetricsHandler
	interval time.Duration
	logger   *zap.SugaredLogger
}

// SendStat отправляет статистику.
func (s *StatSenderService) SendStat(ctx context.Context, ch <-chan []metric.Metrics) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Debug("saveStat done")
			return
		case <-ticker.C:
			var metrics []metric.Metrics

			s.logger.Debug("readCh started")
			if metrics = s.readCh(ch); metrics == nil {
				continue
			}
			s.logger.Debug("readCh finished")

			s.logger.Debug("send started")
			s.send(ctx, metrics)
			s.logger.Debug("send retry finished")
		}
	}
}

func (s *StatSenderService) readCh(ch <-chan []metric.Metrics) []metric.Metrics {
	metrics := make([][]metric.Metrics, 0)
loop:
	for {
		select {
		case m := <-ch:
			if m == nil {
				break loop
			}
			metrics = append(metrics, m)
		default:
			break loop
		}
	}

	if len(metrics) == 0 {
		return nil
	}

	return s.handler.Processing(metrics)
}

func (s *StatSenderService) send(ctx context.Context, metrics []metric.Metrics) {
	retry.Ret(func(_ retry.Policy) bool {
		var err error
		start := time.Now()

		errChan := make(chan error)
		s.sender.SendMetrics(metrics, errChan)

	loop:
		for {
			select {
			case <-ctx.Done():
				s.logger.Debug("send done")
				break loop
			case e, ok := <-errChan:
				if !ok {
					break loop
				}
				if e != nil {
					err = e
					s.logger.Errorf("sendMetrics failed: %s", e)
				}
			}
		}

		s.logger.Debug("SendMetrics time %v", time.Since(start))

		return err == nil
	})
}

// NewStatSenderService конструктор.
func NewStatSenderService(s StatSender, h MetricsHandler, i time.Duration, l *zap.SugaredLogger) *StatSenderService {
	return &StatSenderService{
		sender:   s,
		handler:  h,
		interval: i,
		logger:   l,
	}
}
