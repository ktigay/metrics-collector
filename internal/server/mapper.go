package server

import (
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/server/storage"
)

func mapEntityToMetrics(entity storage.Entity) metric.Metrics {
	m := metric.Metrics{
		ID:    entity.Name,
		MType: string(entity.Type),
	}

	switch t := entity.Value.(type) {
	case int64:
		m.Delta = &t
	case float64:
		m.Value = &t
	}

	return m
}
