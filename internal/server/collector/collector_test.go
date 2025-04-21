package collector

import (
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/server/metric/item"
	"github.com/ktigay/metrics-collector/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricCollector_Save(t *testing.T) {
	type fields struct {
		metrics map[string]*storage.Entity
	}
	type args struct {
		m []item.MetricDTO
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]item.MetricDTO
	}{
		{
			name: "Positive test",
			fields: fields{
				metrics: map[string]*storage.Entity{
					"counter:PollCount": {
						Key:   "counter:PollCount",
						Type:  metric.TypeCounter,
						Name:  metric.PollCount,
						Value: int64(5),
					},
					"gauge:Alloc": {
						Key:   "gauge:Alloc",
						Type:  metric.TypeGauge,
						Name:  string(metric.Alloc),
						Value: 2.0,
					},
				},
			},
			args: args{
				m: []item.MetricDTO{
					{
						Type:  metric.TypeCounter,
						Name:  metric.PollCount,
						Value: "4",
					},
					{
						Type:  metric.TypeGauge,
						Name:  string(metric.Alloc),
						Value: "12.0",
					},
					{
						Type:  metric.TypeGauge,
						Name:  string(metric.BuckHashSys),
						Value: "22.0",
					},
				},
			},
			want: map[string]item.MetricDTO{
				"counter:PollCount": {
					Type:  metric.TypeCounter,
					Name:  metric.PollCount,
					Value: "9",
				},
				"gauge:Alloc": {
					Type:  metric.TypeGauge,
					Name:  string(metric.Alloc),
					Value: "12",
				},
				"gauge:BuckHashSys": {
					Type:  metric.TypeGauge,
					Name:  string(metric.BuckHashSys),
					Value: "22",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewMetricCollector(&storage.MemStorage{
				Metrics: tt.fields.metrics,
			})

			for _, m := range tt.args.m {
				_ = c.Save(m)
			}

			assert.Equal(t, tt.want, c.GetAll())
		})
	}
}
