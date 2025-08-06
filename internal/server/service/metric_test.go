package service

import (
	"context"
	"testing"

	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/server/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMetricCollector_Save(t *testing.T) {
	type fields struct {
		metrics map[string]repository.MetricEntity
	}
	type args struct {
		m []metric.Metrics
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []repository.MetricEntity
	}{
		{
			name: "Positive_test",
			fields: fields{
				metrics: map[string]repository.MetricEntity{
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
				m: []metric.Metrics{
					{
						Type: "counter",
						ID:   "PollCount",
						Delta: func() *int64 {
							x := int64(4)
							return &x
						}(),
					},
					{
						Type: "gauge",
						ID:   "Alloc",
						Value: func() *float64 {
							x := 12.0
							return &x
						}(),
					},
					{
						Type: "gauge",
						ID:   "BuckHashSys",
						Value: func() *float64 {
							x := 22.0
							return &x
						}(),
					},
				},
			},
			want: []repository.MetricEntity{
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
			c := NewMetricCollector(
				&repository.MemMetricRepository{
					Metrics: tt.fields.metrics,
				},
				zap.NewNop().Sugar(),
			)

			for _, m := range tt.args.m {
				_ = c.Save(context.Background(), m)
			}

			var (
				sm  []repository.MetricEntity
				err error
			)
			if sm, err = c.All(context.Background()); err != nil {
				t.Error(err)
			}
			for _, m := range tt.want {
				assert.Contains(t, sm, m)
			}
		})
	}
}
