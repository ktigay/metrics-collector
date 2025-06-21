package collector

import (
	"maps"
	"math/rand/v2"
	"slices"
	"sync/atomic"

	"github.com/ktigay/metrics-collector/internal/metric"
)

// MetricsHandler обработчик собранных метрик.
type MetricsHandler struct {
	counter     atomic.Int64
	randFloatFn func() float64
}

// Processing обрабатывает метрики.
func (s *MetricsHandler) Processing(metrics [][]metric.Metrics) []metric.Metrics {
	merged := s.merge(metrics)
	merged = s.hydrate(merged)

	return merged
}

func (s *MetricsHandler) merge(metrics [][]metric.Metrics) []metric.Metrics {
	merged := make(map[string]metric.Metrics)

	for _, m := range metrics {
		for _, v := range m {
			merged[v.ID] = v
		}
	}
	s.counter.Add(int64(len(metrics)))

	return slices.Collect(maps.Values(merged))
}

func (s *MetricsHandler) hydrate(m []metric.Metrics) []metric.Metrics {
	counter := s.counter.Load()
	rv := s.randFloatFn()
	m = append(m, metric.Metrics{
		ID:    metric.PollCount,
		Type:  string(metric.TypeCounter),
		Delta: &counter,
	}, metric.Metrics{
		ID:    metric.RandomValue,
		Type:  string(metric.TypeGauge),
		Value: &rv,
	})

	return m
}

// NewMetricsHandler конструктор.
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{
		randFloatFn: rand.Float64,
	}
}
