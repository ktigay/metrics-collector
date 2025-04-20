package collector

import (
	"github.com/ktigay/metrics-collector/internal/client/metric/mapper"
	"math/rand"
	"runtime"
)

type RuntimeMetricCollector struct {
	counter uint64
	stat    MetricCollectDTO
}

func NewRuntimeMetricCollector() *RuntimeMetricCollector {
	return &RuntimeMetricCollector{}
}

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

func (c *RuntimeMetricCollector) GetStat() MetricCollectDTO {
	return c.stat
}
