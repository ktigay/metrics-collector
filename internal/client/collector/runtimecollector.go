// Package collector Сборщик метрик.
package collector

import (
	"runtime"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/ktigay/metrics-collector/internal/metric"
)

type (
	readMemStats func(*runtime.MemStats)
	mapper       func(m runtime.MemStats) map[metric.GaugeMetric]float64
)

// RuntimeMetricCollector сборщик метрик.
type RuntimeMetricCollector struct {
	readMem readMemStats
	mapper  mapper
}

// GetStat собирает метрики.
func (c *RuntimeMetricCollector) GetStat() []metric.Metrics {
	var m runtime.MemStats
	c.readMem(&m)

	gaugeMap := c.mapper(m)
	metrics := make([]metric.Metrics, 0, len(gaugeMap))
	typeGauge := string(metric.TypeGauge)
	for k, v := range gaugeMap {
		metrics = append(metrics, metric.Metrics{
			ID:    string(k),
			Type:  typeGauge,
			Value: &v,
		})
	}

	return metrics
}

// NewRuntimeMetricCollector конструктор.
func NewRuntimeMetricCollector() *RuntimeMetricCollector {
	return &RuntimeMetricCollector{
		readMem: runtime.ReadMemStats,
		mapper:  metric.MapGaugeFromMemStats,
	}
}
