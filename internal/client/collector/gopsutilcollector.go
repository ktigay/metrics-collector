package collector

import (
	"time"

	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type (
	virtualMemory func() (*mem.VirtualMemoryStat, error)
	cpuPercent    func(interval time.Duration, percpu bool) ([]float64, error)
)

// GopsUtilCollector структура для метрик из gopsutil.
type GopsUtilCollector struct {
	mem        virtualMemory
	cpuPercent cpuPercent
}

// GetStat собирает метрики.
func (g *GopsUtilCollector) GetStat() []metric.Metrics {
	v, _ := g.mem()
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
				util, _ := g.cpuPercent(0, false)
				return &util[0]
			}(),
		},
	}

	return metrics
}

// NewGopsUtilCollector конструктор.
func NewGopsUtilCollector() *GopsUtilCollector {
	return &GopsUtilCollector{
		mem:        mem.VirtualMemory,
		cpuPercent: cpu.Percent,
	}
}
