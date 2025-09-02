package metric

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveType(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantM   Type
		wantErr bool
	}{
		{
			name: "Positive_test",
			args: args{
				s: "gauge",
			},
			wantM:   TypeGauge,
			wantErr: false,
		},
		{
			name: "Negative test",
			args: args{
				s: "random_string",
			},
			wantM:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotM, err := ResolveType(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantM, gotM)
		})
	}
}

func TestMapGaugeFromMemStats(t *testing.T) {
	type args struct {
		m runtime.MemStats
	}
	tests := []struct {
		want map[GaugeMetric]float64
		name string
		args args
	}{
		{
			name: "Positive_test",
			args: args{
				m: runtime.MemStats{
					Alloc:         1,
					BuckHashSys:   2,
					Frees:         3,
					GCCPUFraction: 4.0123,
					GCSys:         5,
					HeapAlloc:     6,
					HeapIdle:      7,
					HeapInuse:     8,
					HeapObjects:   9,
					HeapReleased:  10,
					HeapSys:       11,
					LastGC:        12,
					Lookups:       13,
					MCacheInuse:   14,
					MCacheSys:     15,
					MSpanInuse:    16,
					MSpanSys:      17,
					Mallocs:       18,
					NextGC:        19,
					NumForcedGC:   20,
					NumGC:         21,
					OtherSys:      22,
					PauseTotalNs:  23,
					StackInuse:    24,
					StackSys:      25,
					Sys:           26,
					TotalAlloc:    27,
				},
			},
			want: map[GaugeMetric]float64{
				Alloc:         1.0,
				BuckHashSys:   2.0,
				Frees:         3.0,
				GCCPUFraction: 4.0123,
				GCSys:         5.0,
				HeapAlloc:     6.0,
				HeapIdle:      7.0,
				HeapInuse:     8.0,
				HeapObjects:   9.0,
				HeapReleased:  10.0,
				HeapSys:       11.0,
				LastGC:        12.0,
				Lookups:       13.0,
				MCacheInuse:   14.0,
				MCacheSys:     15.0,
				MSpanInuse:    16.0,
				MSpanSys:      17.0,
				Mallocs:       18.0,
				NextGC:        19.0,
				NumForcedGC:   20.0,
				NumGC:         21.0,
				OtherSys:      22.0,
				PauseTotalNs:  23.0,
				StackInuse:    24.0,
				StackSys:      25.0,
				Sys:           26.0,
				TotalAlloc:    27.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, MapGaugeFromMemStats(tt.args.m))
		})
	}
}
