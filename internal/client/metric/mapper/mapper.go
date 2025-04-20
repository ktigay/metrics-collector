package mapper

import (
	"github.com/ktigay/metrics-collector/internal/metric"
	"runtime"
)

func MapGaugeFromMemStats(m runtime.MemStats) map[metric.GaugeMetric]float64 {
	return map[metric.GaugeMetric]float64{
		metric.Alloc:         float64(m.Alloc),
		metric.BuckHashSys:   float64(m.BuckHashSys),
		metric.Frees:         float64(m.Frees),
		metric.GCCPUFraction: m.GCCPUFraction,
		metric.GCSys:         float64(m.GCSys),
		metric.HeapAlloc:     float64(m.HeapAlloc),
		metric.HeapIdle:      float64(m.HeapIdle),
		metric.HeapInuse:     float64(m.HeapInuse),
		metric.HeapObjects:   float64(m.HeapObjects),
		metric.HeapReleased:  float64(m.HeapReleased),
		metric.HeapSys:       float64(m.HeapSys),
		metric.LastGC:        float64(m.LastGC),
		metric.Lookups:       float64(m.Lookups),
		metric.MCacheInuse:   float64(m.MCacheInuse),
		metric.MCacheSys:     float64(m.MCacheSys),
		metric.MSpanInuse:    float64(m.MSpanInuse),
		metric.MSpanSys:      float64(m.MSpanSys),
		metric.Mallocs:       float64(m.Mallocs),
		metric.NextGC:        float64(m.NextGC),
		metric.NumForcedGC:   float64(m.NumForcedGC),
		metric.NumGC:         float64(m.NumGC),
		metric.OtherSys:      float64(m.OtherSys),
		metric.PauseTotalNs:  float64(m.PauseTotalNs),
		metric.StackInuse:    float64(m.StackInuse),
		metric.StackSys:      float64(m.StackSys),
		metric.Sys:           float64(m.Sys),
		metric.TotalAlloc:    float64(m.TotalAlloc),
	}
}
