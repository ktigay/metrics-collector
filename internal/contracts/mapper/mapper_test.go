package mapper

import (
	"reflect"
	"testing"

	"github.com/ktigay/metrics-collector/internal/contracts"
	"github.com/ktigay/metrics-collector/internal/metric"
)

func Test_MapFromMetrics(t *testing.T) {
	type args struct {
		m *metric.Metrics
	}
	tests := []struct {
		want *contracts.Metrics
		args args
		name string
	}{
		{
			name: "Counter_metrics",
			args: args{
				m: &metric.Metrics{
					Delta: func() *int64 {
						v := int64(110)
						return &v
					}(),
					ID:   "PollCount",
					Type: "counter",
				},
			},
			want: &contracts.Metrics{
				Delta: 110,
				Id:    "PollCount",
				Type:  "counter",
			},
		},
		{
			name: "Gauge_metrics",
			args: args{
				m: &metric.Metrics{
					Value: func() *float64 {
						v := 110.111
						return &v
					}(),
					ID:   "Mallocs",
					Type: "gauge",
				},
			},
			want: &contracts.Metrics{
				Value: 110.111,
				Id:    "Mallocs",
				Type:  "gauge",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapFromMetrics(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapFromMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_MapToMetrics(t *testing.T) {
	type args struct {
		m *contracts.Metrics
	}
	tests := []struct {
		name string
		args args
		want metric.Metrics
	}{
		{
			name: "Counter_metrics",
			args: args{
				m: &contracts.Metrics{
					Delta: 110,
					Id:    "PollCount",
					Type:  "counter",
				},
			},
			want: metric.Metrics{
				Delta: func() *int64 {
					v := int64(110)
					return &v
				}(),
				ID:   "PollCount",
				Type: "counter",
			},
		},
		{
			name: "Gauge_metrics",
			args: args{
				m: &contracts.Metrics{
					Value: 110.111,
					Id:    "Mallocs",
					Type:  "gauge",
				},
			},
			want: metric.Metrics{
				Value: func() *float64 {
					v := 110.111
					return &v
				}(),
				ID:   "Mallocs",
				Type: "gauge",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapToMetrics(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
