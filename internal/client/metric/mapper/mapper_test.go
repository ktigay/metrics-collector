package mapper

import (
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestMapGaugeFromMemStats(t *testing.T) {
	type args struct {
		m runtime.MemStats
	}
	tests := []struct {
		name string
		args args
		want map[metric.GaugeMetric]float64
	}{
		{
			name: "Positive test",
			args: args{
				m: runtime.MemStats{
					Alloc:         1,
					BuckHashSys:   2,
					Frees:         3,
					GCCPUFraction: 4.0123,
					GCSys:         5,
					HeapAlloc:     6,
					HeapIdle:      7,
					HeapInuse:     8,
					HeapObjects:   9,
					HeapReleased:  10,
					HeapSys:       11,
					LastGC:        12,
					Lookups:       13,
					MCacheInuse:   14,
					MCacheSys:     15,
					MSpanInuse:    16,
					MSpanSys:      17,
					Mallocs:       18,
					NextGC:        19,
					NumForcedGC:   20,
					NumGC:         21,
					OtherSys:      22,
					PauseTotalNs:  23,
					StackInuse:    24,
					StackSys:      25,
					Sys:           26,
					TotalAlloc:    27,
				},
			},
			want: map[metric.GaugeMetric]float64{
				metric.Alloc:         1.0,
				metric.BuckHashSys:   2.0,
				metric.Frees:         3.0,
				metric.GCCPUFraction: 4.0123,
				metric.GCSys:         5.0,
				metric.HeapAlloc:     6.0,
				metric.HeapIdle:      7.0,
				metric.HeapInuse:     8.0,
				metric.HeapObjects:   9.0,
				metric.HeapReleased:  10.0,
				metric.HeapSys:       11.0,
				metric.LastGC:        12.0,
				metric.Lookups:       13.0,
				metric.MCacheInuse:   14.0,
				metric.MCacheSys:     15.0,
				metric.MSpanInuse:    16.0,
				metric.MSpanSys:      17.0,
				metric.Mallocs:       18.0,
				metric.NextGC:        19.0,
				metric.NumForcedGC:   20.0,
				metric.NumGC:         21.0,
				metric.OtherSys:      22.0,
				metric.PauseTotalNs:  23.0,
				metric.StackInuse:    24.0,
				metric.StackSys:      25.0,
				metric.Sys:           26.0,
				metric.TotalAlloc:    27.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, MapGaugeFromMemStats(tt.args.m))
		})
	}
}
