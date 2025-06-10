// Package collector Сборщик метрик.
package collector

import (
	"context"
	"runtime"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/ktigay/metrics-collector/internal/metric"
)

//go:generate mockgen -destination=./mocks/mock_storage.go -package=mocks github.com/ktigay/metrics-collector/internal/client/collector StorageInterface
type StorageInterface interface {
	Save(metrics []metric.Metrics)
}

// RuntimeMetricCollector сборщик метрик.
type RuntimeMetricCollector struct {
	storage StorageInterface
}

// NewRuntimeMetricCollector конструктор.
func NewRuntimeMetricCollector(storage *Storage) *RuntimeMetricCollector {
	return &RuntimeMetricCollector{
		storage: storage,
	}
}

// PollStat собирает метрики.
func (c *RuntimeMetricCollector) PollStat(_ context.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	gaugeMap := metric.MapGaugeFromMemStats(m)
	metrics := make([]metric.Metrics, 0, len(gaugeMap))
	typeGauge := string(metric.TypeGauge)
	for k, v := range gaugeMap {
		metrics = append(metrics, metric.Metrics{
			ID:    string(k),
			MType: typeGauge,
			Value: &v,
		})
	}

	c.storage.Save(metrics)
}
