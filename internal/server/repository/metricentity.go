package repository

import (
	"github.com/ktigay/metrics-collector/internal/metric"
)

// MetricEntity сущность для сохранения в repository.
type MetricEntity struct {
	Key   string      `json:"key"`
	Type  metric.Type `json:"type"`
	Name  string      `json:"name"`
	Delta int64       `json:"delta"`
	Value float64     `json:"value"`
}

// ValueByType возвращает значение в зависимости от типа.
func (e *MetricEntity) ValueByType() any {
	switch e.Type {
	case metric.TypeCounter:
		return e.Delta
	case metric.TypeGauge:
		return e.Value
	}
	return nil
}

// ToMetrics мап сущности в дто.
func (e *MetricEntity) ToMetrics() metric.Metrics {
	m := metric.Metrics{
		ID:    e.Name,
		MType: string(e.Type),
	}

	switch e.Type {
	case metric.TypeCounter:
		m.Delta = &e.Delta
	case metric.TypeGauge:
		m.Value = &e.Value
	}

	return m
}
