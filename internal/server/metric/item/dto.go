package item

import "github.com/ktigay/metrics-collector/internal/metric"

// MetricDTO - объект метрики.
type MetricDTO struct {
	Type       metric.Type
	Name       string
	IntValue   int64
	FloatValue float64
}

// GetKey - возвращает уникальный ключ метрики.
func (m *MetricDTO) GetKey() string {
	return string(m.Type) + ":" + m.Name
}
