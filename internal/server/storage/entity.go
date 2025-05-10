package storage

import (
	"github.com/ktigay/metrics-collector/internal/metric"
)

// Entity - сущность для сохранения в storage.
type Entity struct {
	Key   string      `json:"key"`
	Type  metric.Type `json:"type"`
	Name  string      `json:"name"`
	Delta int64       `json:"delta"`
	Value float64     `json:"value"`
}

// ValueByType возвращает значение в зависимости от типа.
func (e *Entity) ValueByType() any {
	switch e.Type {
	case metric.TypeCounter:
		return e.Delta
	case metric.TypeGauge:
		return e.Value
	}
	return nil
}

// MapEntityToMetrics мап сущности в дто.
func MapEntityToMetrics(entity Entity) metric.Metrics {
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
