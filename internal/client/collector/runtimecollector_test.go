package collector

import (
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRuntimeMetricCollector_PollStat(t *testing.T) {
	type fields struct {
		counter uint64
		stat    MetricCollectDTO
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{
			name: "Positive test",
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

			c.PollStat()

			assert.Equal(t, tt.want, c.counter)
			assert.Equal(t, tt.want, c.GetStat().Counter)
		})
	}
}
