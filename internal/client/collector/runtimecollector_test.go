package collector

import (
	"context"
	"testing"

	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/stretchr/testify/assert"
)

func TestRuntimeMetricCollector_PollStat(t *testing.T) {
	type fields struct {
		counter int64
		stat    MetricCollectDTO
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name: "Positive_test",
			fields: fields{
				stat: MetricCollectDTO{
					MemStats: map[metric.GaugeMetric]float64{},
					Rand:     1222.222,
				},
				counter: 12,
			},
			want: 13,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RuntimeMetricCollector{
				counter: tt.fields.counter,
				stat:    tt.fields.stat,
			}

			c.PollStat(context.TODO())

			assert.Equal(t, tt.want, c.counter)
			assert.Equal(t, tt.want, c.GetStat().Counter)
		})
	}
}
