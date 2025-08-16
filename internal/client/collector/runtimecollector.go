// Package collector Сборщик метрик.
package collector

import (
	"runtime"

	_ "github.com/golang/mock/mockgen/model"

	"github.com/ktigay/metrics-collector/internal/metric"
)

type (
	readMemStatsFn func(*runtime.MemStats)
	mapperFn       func(m runtime.MemStats) map[metric.GaugeMetric]float64
)

// RuntimeMetricCollector сборщик метрик.
type RuntimeMetricCollector struct {
	readMemFn readMemStatsFn
	mapperFn  mapperFn
}

// GetStat собирает метрики.
func (c *RuntimeMetricCollector) GetStat() ([]metric.Metrics, error) {
	var m runtime.MemStats
	c.readMemFn(&m)

	gaugeMap := c.mapperFn(m)
	metrics := make([]metric.Metrics, 0, len(gaugeMap))
	typeGauge := string(metric.TypeGauge)
	for k, v := range gaugeMap {
		metrics = append(metrics, metric.Metrics{
			ID:    string(k),
			Type:  typeGauge,
			Value: &v,
		})
	}

	return metrics, nil
}

// NewRuntimeMetricCollector конструктор.
func NewRuntimeMetricCollector() *RuntimeMetricCollector {
	return &RuntimeMetricCollector{
		readMemFn: runtime.ReadMemStats,
		mapperFn:  metric.MapGaugeFromMemStats,
	}
}
