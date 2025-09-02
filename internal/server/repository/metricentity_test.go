package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ktigay/metrics-collector/internal/metric"
)

func TestEntity_ToMetrics(t *testing.T) {
	type fields struct {
		Key   string
		Type  metric.Type
		Name  string
		Delta int64
		Value float64
	}
	tests := []struct {
		want   metric.Metrics
		name   string
		fields fields
	}{
		{
			name: "TestEntity_ToMetrics_gauge",
			fields: fields{
				Key:   "gauge:Mallocs",
				Type:  metric.TypeGauge,
				Name:  "Mallocs",
				Delta: 0,
				Value: 12.33345,
			},
			want: metric.Metrics{
				ID:    "Mallocs",
				Type:  "gauge",
				Delta: nil,
				Value: func() *float64 {
					v := 12.33345
					return &v
				}(),
			},
		},
		{
			name: "TestEntity_ToMetrics_counter",
			fields: fields{
				Key:   "counter:PollCount",
				Type:  metric.TypeCounter,
				Name:  "PollCount",
				Delta: 120,
				Value: .0,
			},
			want: metric.Metrics{
				ID:   "PollCount",
				Type: "counter",
				Delta: func() *int64 {
					v := int64(120)
					return &v
				}(),
				Value: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &MetricEntity{
				Key:   tt.fields.Key,
				Type:  tt.fields.Type,
				Name:  tt.fields.Name,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			assert.Equalf(t, tt.want, e.ToMetrics(), "ToMetrics()")
		})
	}
}

func TestEntity_ValueByType(t *testing.T) {
	type fields struct {
		Key   string
		Type  metric.Type
		Name  string
		Delta int64
		Value float64
	}
	tests := []struct {
		want   any
		name   string
		fields fields
	}{
		{
			name: "TestEntity_ValueByType_gauge",
			fields: fields{
				Key:   "gauge:Mallocs",
				Type:  metric.TypeGauge,
				Name:  "Mallocs",
				Delta: 0,
				Value: 12.33345,
			},
			want: 12.33345,
		},
		{
			name: "TestEntity_ValueByType_counter",
			fields: fields{
				Key:   "counter:PollCount",
				Type:  metric.TypeCounter,
				Name:  "PollCount",
				Delta: 120,
				Value: .0,
			},
			want: int64(120),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &MetricEntity{
				Key:   tt.fields.Key,
				Type:  tt.fields.Type,
				Name:  tt.fields.Name,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			assert.Equalf(t, tt.want, e.ValueByType(), "ValueByType()")
		})
	}
}
