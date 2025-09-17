package handler

import (
	"context"

	"github.com/ktigay/metrics-collector/internal/metric"
)

// Collector Интерфейс сборщика статистики.
//
//go:generate mockgen -destination=./mocks/mock_collector.go -package=mocks github.com/ktigay/metrics-collector/internal/server/handler Collector
type Collector interface {
	Save(ctx context.Context, mt metric.Metrics) (*metric.Metrics, error)
	All(ctx context.Context) (*[]metric.Metrics, error)
	Find(ctx context.Context, t, n string) (*metric.Metrics, error)
	Remove(ctx context.Context, t, n string) error
	SaveAll(ctx context.Context, mt []metric.Metrics) error
}
