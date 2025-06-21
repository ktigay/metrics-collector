package service

import (
	"context"

	"github.com/ktigay/metrics-collector/internal/metric"
	"go.uber.org/zap"
)

// StorageInterface интерфейс репозитория.
//
//go:generate mockgen -destination=./mocks/mock_storage.go -package=mocks github.com/ktigay/metrics-collector/internal/client/service StorageInterface
type StorageInterface interface {
	Save(metrics []metric.Metrics)
}

// StatSaverService структура для сборки статистики.
type StatSaverService struct {
	storage StorageInterface
	logger  *zap.SugaredLogger
}

// PushStat собирает статистику.
func (s *StatSaverService) PushStat(ctx context.Context, ch <-chan []metric.Metrics) {
	for {
		select {
		case <-ctx.Done():
			s.logger.Debug("saveStat done")
			return
		case m, ok := <-ch:
			if !ok {
				s.logger.Debug("saveStat channel closed")
				return
			}
			s.storage.Save(m)
		}
	}
}

// NewStatSaverService конструктор.
func NewStatSaverService(storage StorageInterface, logger *zap.SugaredLogger) *StatSaverService {
	return &StatSaverService{
		storage: storage,
		logger:  logger,
	}
}
