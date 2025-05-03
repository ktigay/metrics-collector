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

	switch entity.Type {
	case metric.TypeCounter:
		m.Delta = &entity.Delta
	case metric.TypeGauge:
		m.Value = &entity.Value
	}

	return m
}
