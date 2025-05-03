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

// GetKey - возвращает уникальный ключ метрики.
func (e *Entity) GetKey() string {
	return metric.Key(string(e.Type), e.Name)
}

// GetValue - возвращает значение для типа gauge.
func (e *Entity) GetValue() any {
	return e.Value
}

// GetDelta возвращает значение для типа counter.
func (e *Entity) GetDelta() int64 {
	return e.Delta
}

// GetValueByType возвращает значение в зависимости от типа.
func (e *Entity) GetValueByType() any {
	switch e.Type {
	case metric.TypeCounter:
		return e.Delta
	case metric.TypeGauge:
		return e.Value
	}
	return nil
}
