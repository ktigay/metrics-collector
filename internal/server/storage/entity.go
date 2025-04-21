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
	Value interface{}
}

// GetValueAsString - возвращает значение как строку.
func (s *Entity) GetValueAsString() string {
	if s == nil {
		return ""
	}
	switch s.Value.(type) {
	case int64:
		return fmt.Sprintf("%d", s.Value)
	case float64:
		return fmt.Sprintf("%g", s.Value)
	}
	return ""
}
