package metric

import "fmt"

// Type - тип метрики.
type Type string

const (
	TypeGauge   Type = "gauge"
	TypeCounter Type = "counter"
)

// ResolveType - получает из строки тип.
func ResolveType(s string) (m Type, err error) {
	switch s {
	case string(TypeGauge):
		return TypeGauge, nil
	case string(TypeCounter):
		return TypeCounter, nil
	default:
		return m, fmt.Errorf("unknown metric type: %s", s)
	}
}
