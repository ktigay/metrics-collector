package collector

import (
	"runtime"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ktigay/metrics-collector/internal/metric"
)

func TestRuntimeMetricCollector_GetStat(t *testing.T) {
	type fields struct {
		readMem readMemStatsFn
		mapper  mapperFn
	}
	tests := []struct {
		name   string
		fields fields
		want   []metric.Metrics
	}{
		{
			name: "Positive_test_GetStat",
			fields: fields{
				readMem: func(stats *runtime.MemStats) {
					stats.Alloc = 100
					stats.TotalAlloc = 200
				},
				mapper: func(m runtime.MemStats) map[metric.GaugeMetric]float64 {
					return map[metric.GaugeMetric]float64{
						metric.Alloc:      float64(m.Alloc),
						metric.TotalAlloc: float64(m.TotalAlloc),
					}
				},
			},
			want: []metric.Metrics{
				{
					ID:   "Alloc",
					Type: "gauge",
					Value: func() *float64 {
						v := 100.0
						return &v
					}(),
				},
				{
					ID:   "TotalAlloc",
					Type: "gauge",
					Value: func() *float64 {
						v := 200.0
						return &v
					}(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RuntimeMetricCollector{
				readMemFn: tt.fields.readMem,
				mapperFn:  tt.fields.mapper,
			}
			got := c.GetStat()
			sort.Slice(got, func(i, j int) bool { return got[i].ID < got[j].ID })
			sort.Slice(tt.want, func(i, j int) bool { return tt.want[i].ID < tt.want[j].ID })

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("GetStat() = %v, want %v, diff %v", got, tt.want, diff)
			}
		})
	}
}
