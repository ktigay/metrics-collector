package item

import (
	"github.com/ktigay/metrics-collector/internal/metric"
)

// MetricDTO - объект метрики.
type MetricDTO struct {
	Type  metric.Type
	Name  string
	Value string
}

// GetKey - возвращает уникальный ключ метрики.
func (m *MetricDTO) GetKey() string {
	return metric.GetKey(string(m.Type), m.Name)
}

// GetValue - возвращает значение.
func (m *MetricDTO) GetValue() string {
	return m.Value
}
