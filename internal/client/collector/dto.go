package collector

import "github.com/ktigay/metrics-collector/internal/metric"

// MetricCollectDTO - DTO.
type MetricCollectDTO struct {
	MemStats map[metric.GaugeMetric]float64
	Counter  uint64
	Rand     float64
}
