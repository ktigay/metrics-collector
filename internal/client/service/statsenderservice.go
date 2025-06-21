package service

import (
	"context"
	"time"

	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/retry"
	"go.uber.org/zap"
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
			metrics := make([][]metric.Metrics, 0)
		loop:
			for {
				select {
				case m := <-ch:
					metrics = append(metrics, m)
				default:
					break loop
				}
			}

			s.send(s.handler.Processing(metrics))
		}
	}
}

func (s *StatSenderService) send(metrics []metric.Metrics) {
	s.logger.Debug("send started")

	retry.Ret(func(_ retry.Policy) bool {
		var err error
		start := time.Now()

		errChan := make(chan error)
		s.sender.SendMetrics(metrics, errChan)

		for e := range errChan {
			if e != nil {
				err = e
				s.logger.Errorf("sendMetrics failed: %s", e)
			}
		}

		s.logger.Debug("SendMetrics time %v", time.Since(start))

		return err == nil
	})

	s.logger.Debug("send retry finished")
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
