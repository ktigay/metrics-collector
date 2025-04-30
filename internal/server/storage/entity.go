package storage

import (
	"fmt"
	"github.com/ktigay/metrics-collector/internal/metric"
)

// Entity - сущность для сохранения в storage.
type Entity struct {
	Key   string
	Type  metric.Type
	Name  string
	Value any
}

// GetKey - возвращает уникальный ключ метрики.
func (e *Entity) GetKey() string {
	return metric.GetKey(string(e.Type), e.Name)
}

// GetValue - возвращает значение.
func (e *Entity) GetValue() any {
	return e.Value
}

// ValueAsString - возвращает значение как строку.
func (e *Entity) ValueAsString() string {
	if e == nil {
		return ""
	}
	switch e.Value.(type) {
	case int64:
		return fmt.Sprintf("%d", e.Value)
	case float64:
		return fmt.Sprintf("%g", e.Value)
	}
	return ""
}
