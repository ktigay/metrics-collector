package service

import (
	"context"
	"time"

	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/retry"
	"go.uber.org/zap"
)

// StatSender отправка метрик.
//
//go:generate mockgen -destination=./mocks/mock_sender.go -package=mocks github.com/ktigay/metrics-collector/internal/client/service StatSender
type StatSender interface {
	SendMetrics([]metric.Metrics, chan<- error)
}

// StatSenderService провайдер статистики.
type StatSenderService struct {
	sender StatSender
	logger *zap.SugaredLogger
}

// SendStat отправляет статистику.
func (s *StatSenderService) SendStat(ctx context.Context, ch <-chan []metric.Metrics) {
	for {
		select {
		case <-ctx.Done():
			s.logger.Debug("sendStat done")
			return
		case metrics, ok := <-ch:
			if !ok {
				s.logger.Debug("sendStat channel closed")
				return
			}

			s.logger.Debug("sendStat started")

			retry.Ret(func(_ retry.Policy) bool {
				select {
				case <-ctx.Done():
					s.logger.Debug("retry.Ret stopped")
					return true
				default:
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
				}
			})

			s.logger.Debug("sendStat retry finished")
		}
	}
}

// NewStatSenderService конструктор.
func NewStatSenderService(s StatSender, logger *zap.SugaredLogger) *StatSenderService {
	return &StatSenderService{
		sender: s,
		logger: logger,
	}
}
