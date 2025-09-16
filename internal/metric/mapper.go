package metric

import "github.com/ktigay/metrics-collector/internal/contracts"

// MapToMetrics из grpc в структуру.
func MapToMetrics(m *contracts.Metrics) Metrics {
	return Metrics{
		Delta: func() *int64 {
			delta := m.GetDelta()
			if delta == 0 {
				return nil
			}
			return &delta
		}(),
		Value: func() *float64 {
			val := m.GetValue()
			if val == .0 {
				return nil
			}
			return &val
		}(),
		ID:   m.GetId(),
		Type: m.GetType(),
	}
}

// MapFromMetrics из структуры в grpc.
func MapFromMetrics(m *Metrics) *contracts.Metrics {
	return &contracts.Metrics{
		Delta: m.GetDelta(),
		Value: m.GetValue(),
		Type:  m.Type,
		Id:    m.ID,
	}
}
