package collector

import (
	"time"

	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type (
	virtualMemoryFn func() (*mem.VirtualMemoryStat, error)
	cpuPercentFn    func(interval time.Duration, percpu bool) ([]float64, error)
)

// GopsUtilCollector структура для метрик из gopsutil.
type GopsUtilCollector struct {
	memFn        virtualMemoryFn
	cpuPercentFn cpuPercentFn
}

// GetStat собирает метрики.
func (g *GopsUtilCollector) GetStat() []metric.Metrics {
	v, _ := g.memFn()
	typeGauge := string(metric.TypeGauge)

	metrics := []metric.Metrics{
		{
			ID:   string(metric.TotalMemory),
			Type: typeGauge,
			Value: func() *float64 {
				total := float64(v.Total)
				return &total
			}(),
		},
		{
			ID:   string(metric.FreeMemory),
			Type: typeGauge,
			Value: func() *float64 {
				free := float64(v.Free)
				return &free
			}(),
		},
		{
			ID:   string(metric.CPUutilization1),
			Type: typeGauge,
			Value: func() *float64 {
				util, _ := g.cpuPercentFn(0, false)
				return &util[0]
			}(),
		},
	}

	return metrics
}

// NewGopsUtilCollector конструктор.
func NewGopsUtilCollector() *GopsUtilCollector {
	return &GopsUtilCollector{
		memFn:        mem.VirtualMemory,
		cpuPercentFn: cpu.Percent,
	}
}
