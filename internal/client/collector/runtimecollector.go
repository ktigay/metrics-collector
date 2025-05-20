package collector

import (
	"math/rand/v2"
	"runtime"

	"github.com/ktigay/metrics-collector/internal/metric"
)

// RuntimeMetricCollector - сборщик метрик.
type RuntimeMetricCollector struct {
	counter int64
	stat    MetricCollectDTO
}

// MetricCollectDTO - DTO.
type MetricCollectDTO struct {
	MemStats map[metric.GaugeMetric]float64
	Counter  int64
	Rand     float64
}

// NewRuntimeMetricCollector - конструктор.
func NewRuntimeMetricCollector() *RuntimeMetricCollector {
	return &RuntimeMetricCollector{}
}

// PollStat - собирает метрики.
func (c *RuntimeMetricCollector) PollStat() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	c.counter++

	c.stat = MetricCollectDTO{
		MemStats: metric.MapGaugeFromMemStats(m),
		Counter:  c.counter,
		Rand:     rand.Float64(),
	}
}

// GetStat - возвращает метрики.
func (c *RuntimeMetricCollector) GetStat() MetricCollectDTO {
	return c.stat
}
