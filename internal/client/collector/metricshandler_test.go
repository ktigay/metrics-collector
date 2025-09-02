package collector

import (
	"sort"
	"sync/atomic"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/ktigay/metrics-collector/internal/metric"
)

func TestMetricsHandler_Processing(t *testing.T) {
	type fields struct {
		randFloatFn func() float64
		counter     int64
	}
	type args struct {
		metrics [][]metric.Metrics
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []metric.Metrics
	}{
		{
			name: "Positive_test_Processing",
			fields: fields{
				counter: 100,
				randFloatFn: func() float64 {
					return 50.5
				},
			},
			args: args{
				metrics: [][]metric.Metrics{
					{
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
					{
						{
							ID:   "Alloc",
							Type: "gauge",
							Value: func() *float64 {
								v := 123.321
								return &v
							}(),
						},
						{
							ID:   "HeapAlloc",
							Type: "gauge",
							Value: func() *float64 {
								v := 150.123
								return &v
							}(),
						},
					},
				},
			},
			want: []metric.Metrics{
				{
					ID:   "Alloc",
					Type: "gauge",
					Value: func() *float64 {
						v := 123.321
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
				{
					ID:   "HeapAlloc",
					Type: "gauge",
					Value: func() *float64 {
						v := 150.123
						return &v
					}(),
				},
				{
					ID:   "RandomValue",
					Type: "gauge",
					Value: func() *float64 {
						v := 50.5
						return &v
					}(),
				},
				{
					ID:   "PollCount",
					Type: "counter",
					Delta: func() *int64 {
						v := int64(102)
						return &v
					}(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsHandler{
				counter:     atomic.Int64{},
				randFloatFn: tt.fields.randFloatFn,
			}
			s.counter.Add(tt.fields.counter)

			got := s.Processing(tt.args.metrics)
			sort.Slice(got, func(i, j int) bool { return got[i].ID < got[j].ID })
			sort.Slice(tt.want, func(i, j int) bool { return tt.want[i].ID < tt.want[j].ID })

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Processing() = %v, want %v, diff %v", got, tt.want, diff)
			}
		})
	}
}
