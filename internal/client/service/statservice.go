package service

import (
	"context"
	"time"

	"github.com/ktigay/metrics-collector/internal/client/collector"
	"github.com/ktigay/metrics-collector/internal/retry"
	"go.uber.org/zap"
)

// StatSender отправка метрик.
type StatSender interface {
	SendMetrics(c collector.MetricCollectDTO) error
}

// StatGetter получение метрик.
type StatGetter interface {
	GetStat() collector.MetricCollectDTO
}

// StatService провайдер статистики.
type StatService struct {
	getter       StatGetter
	sender       StatSender
	sendInterval time.Duration
	logger       *zap.SugaredLogger
}

// SendStat отправляет статистику.
func (s *StatService) SendStat(ctx context.Context) {
	ticker := time.NewTicker(s.sendInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			s.logger.Debug("sendStat stopped")
			return
		case <-ticker.C:
			s.logger.Debug("sendStat started")

			retry.Ret(func(_ retry.Policy) bool {
				select {
				case <-ctx.Done():
					s.logger.Debug("retry.Ret stopped")
					return true
				default:
					err := s.sender.SendMetrics(s.getter.GetStat())
					if err != nil {
						s.logger.Errorf("sendStat err: %v", err)
					}
					return err == nil
				}
			})

			s.logger.Debug("sendStat finished")
		}
	}
}

// NewStatService конструктор.
func NewStatService(cl StatGetter, s StatSender, sendInterval time.Duration, logger *zap.SugaredLogger) *StatService {
	return &StatService{
		getter:       cl,
		sender:       s,
		sendInterval: sendInterval,
		logger:       logger,
	}
}
