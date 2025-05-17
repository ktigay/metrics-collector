package storage

import (
	"testing"

	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/stretchr/testify/assert"
)

func TestMetricCollector_Save(t *testing.T) {
	type fields struct {
		metrics map[string]Entity
	}
	type args struct {
		m []struct {
			Type  metric.Type
			Name  string
			Value string
		}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Entity
	}{
		{
			name: "Positive_test",
			fields: fields{
				metrics: map[string]Entity{
					"counter:PollCount": {
						Key:   "counter:PollCount",
						Type:  metric.TypeCounter,
						Name:  metric.PollCount,
						Delta: int64(5),
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
				m: []struct {
					Type  metric.Type
					Name  string
					Value string
				}{
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
			want: []Entity{
				{
					Key:   "counter:PollCount",
					Type:  metric.TypeCounter,
					Name:  metric.PollCount,
					Delta: int64(9),
				},
				{
					Key:   "gauge:Alloc",
					Type:  metric.TypeGauge,
					Name:  string(metric.Alloc),
					Value: 12.000,
				},
				{
					Key:   "gauge:BuckHashSys",
					Type:  metric.TypeGauge,
					Name:  string(metric.BuckHashSys),
					Value: 22.000,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewMetricCollector(&MemStorage{
				Metrics: tt.fields.metrics,
			})

			for _, m := range tt.args.m {
				_ = c.Save(m.Type, m.Name, m.Value)
			}

			a := c.GetAll()
			for _, m := range tt.want {
				assert.Contains(t, a, m)
			}
		})
	}
}
