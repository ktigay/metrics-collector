package collector

import (
	"context"

	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// GopsUtilCollector структура для метрик из gopsutil.
type GopsUtilCollector struct {
	storage StorageInterface
}

// PollStat собирает метрики.
func (g *GopsUtilCollector) PollStat(_ context.Context) {
	v, _ := mem.VirtualMemory()
	typeGauge := string(metric.TypeGauge)

	metrics := []metric.Metrics{
		{
			ID:    string(metric.TotalMemory),
			MType: typeGauge,
			Value: func() *float64 {
				total := float64(v.Total)
				return &total
			}(),
		},
		{
			ID:    string(metric.FreeMemory),
			MType: typeGauge,
			Value: func() *float64 {
				free := float64(v.Free)
				return &free
			}(),
		},
		{
			ID:    string(metric.CPUutilization1),
			MType: typeGauge,
			Value: func() *float64 {
				util, _ := cpu.Percent(0, false)
				return &util[0]
			}(),
		},
	}

	g.storage.Save(metrics)
}

// NewGopsUtilCollector конструктор.
func NewGopsUtilCollector(storage *Storage) *GopsUtilCollector {
	return &GopsUtilCollector{
		storage: storage,
	}
}
