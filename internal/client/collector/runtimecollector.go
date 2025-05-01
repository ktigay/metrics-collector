package collector

import (
	"github.com/ktigay/metrics-collector/internal/client/metric/mapper"
	"math/rand"
	"runtime"
)

// RuntimeMetricCollector - сборщик метрик.
type RuntimeMetricCollector struct {
	counter int64
	stat    MetricCollectDTO
}

// NewRuntimeMetricCollector - конструктор.
func NewRuntimeMetricCollector() *RuntimeMetricCollector {
	return &RuntimeMetricCollector{}
}

// PollStat - собирает метрики.
func (c *RuntimeMetricCollector) PollStat() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	c.counter += 1

	c.stat = MetricCollectDTO{
		MemStats: mapper.MapGaugeFromMemStats(m),
		Counter:  c.counter,
		Rand:     rand.Float64(),
	}
}

// GetStat - возвращает метрики.
func (c *RuntimeMetricCollector) GetStat() MetricCollectDTO {
	return c.stat
}
