package collector

import (
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/shirou/gopsutil/v4/mem"
)

func TestGopsUtilCollector_GetStat(t *testing.T) {
	type fields struct {
		memFn        virtualMemoryFn
		cpuPercentFn cpuPercentFn
	}
	tests := []struct {
		name   string
		fields fields
		want   []metric.Metrics
	}{
		{
			name: "Positive_test_GetStat",
			fields: fields{
				memFn: func() (*mem.VirtualMemoryStat, error) {
					stat := &mem.VirtualMemoryStat{
						Total: 100,
						Free:  200,
					}
					return stat, nil
				},
				cpuPercentFn: func(interval time.Duration, percpu bool) ([]float64, error) {
					return []float64{10, 20}, nil
				},
			},
			want: []metric.Metrics{
				{
					ID:   "TotalMemory",
					Type: "gauge",
					Value: func() *float64 {
						v := 100.0
						return &v
					}(),
				},
				{
					ID:   "FreeMemory",
					Type: "gauge",
					Value: func() *float64 {
						v := 200.0
						return &v
					}(),
				},
				{
					ID:   "CPUutilization1",
					Type: "gauge",
					Value: func() *float64 {
						v := 10.0
						return &v
					}(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GopsUtilCollector{
				memFn:        tt.fields.memFn,
				cpuPercentFn: tt.fields.cpuPercentFn,
			}
			got := g.GetStat()
			sort.Slice(got, func(i, j int) bool { return got[i].ID < got[j].ID })
			sort.Slice(tt.want, func(i, j int) bool { return tt.want[i].ID < tt.want[j].ID })

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("GetStat() = %v, want %v, diff %v", got, tt.want, diff)
			}
		})
	}
}
