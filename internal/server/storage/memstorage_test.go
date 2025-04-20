package storage

import (
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/server/metric/item"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemStorage_Save(t *testing.T) {
	type fields struct {
		metrics map[string]*MemStorageEntity
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
				metrics: map[string]*MemStorageEntity{
					"counter:PollCount": {
						Key:        "counter:PollCount",
						Type:       metric.TypeCounter,
						Name:       metric.PollCount,
						IntValue:   5,
						FloatValue: 0,
					},
					"gauge:Alloc": {
						Key:        "gauge:Alloc",
						Type:       metric.TypeGauge,
						Name:       string(metric.Alloc),
						IntValue:   0,
						FloatValue: 2.0,
					},
				},
			},
			args: args{
				m: []item.MetricDTO{
					{
						Type:       metric.TypeCounter,
						Name:       metric.PollCount,
						IntValue:   4,
						FloatValue: 0,
					},
					{
						Type:       metric.TypeGauge,
						Name:       string(metric.Alloc),
						IntValue:   0,
						FloatValue: 12.0,
					},
					{
						Type:       metric.TypeGauge,
						Name:       string(metric.BuckHashSys),
						IntValue:   0,
						FloatValue: 22.0,
					},
				},
			},
			want: map[string]item.MetricDTO{
				"counter:PollCount": {
					Type:       metric.TypeCounter,
					Name:       metric.PollCount,
					IntValue:   9,
					FloatValue: 0,
				},
				"gauge:Alloc": {
					Type:       metric.TypeGauge,
					Name:       string(metric.Alloc),
					IntValue:   0,
					FloatValue: 12.0,
				},
				"gauge:BuckHashSys": {
					Type:       metric.TypeGauge,
					Name:       string(metric.BuckHashSys),
					IntValue:   0,
					FloatValue: 22.0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				metrics: tt.fields.metrics,
			}

			for _, m := range tt.args.m {
				_ = s.Save(m)
			}

			assert.Equal(t, tt.want, s.GetAll())
		})
	}
}
