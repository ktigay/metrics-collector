package client

import (
	"github.com/ktigay/metrics-collector/internal/client/collector"
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestMetricSender_sendCounter(t *testing.T) {
	type args struct {
		c collector.MetricCollectDTO
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Positive_test",
			args: args{
				c: collector.MetricCollectDTO{
					Counter: 15,
				},
			},
			want: "/update/counter/PollCount/15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := http.NewServeMux()
			router.HandleFunc("/update/{type}/{name}/{value}", func(writer http.ResponseWriter, request *http.Request) {
				require.Equal(t, http.MethodPost, request.Method)
				assert.Equal(t, tt.want, request.RequestURI)
			})

			svr := httptest.NewServer(router)
			defer svr.Close()

			c := NewMetricHandler(svr.URL)
			_ = c.sendCounter(tt.args.c)
		})
	}
}

func TestMetricSender_sendGaugeMetrics(t *testing.T) {
	type args struct {
		c collector.MetricCollectDTO
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Positive_test",
			args: args{
				c: collector.MetricCollectDTO{
					MemStats: map[metric.GaugeMetric]float64{
						metric.Alloc: 12.345,
					},
				},
			},
			want: "/update/gauge/Alloc/12.345",
		},
		{
			name: "Positive_test_#2",
			args: args{
				c: collector.MetricCollectDTO{
					MemStats: map[metric.GaugeMetric]float64{
						metric.HeapReleased: 112.3245,
					},
				},
			},
			want: "/update/gauge/HeapReleased/112.3245",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := http.NewServeMux()
			router.HandleFunc("/update/{type}/{name}/{value}", func(writer http.ResponseWriter, request *http.Request) {
				require.Equal(t, http.MethodPost, request.Method)
				assert.Equal(t, tt.want, request.RequestURI)
			})

			svr := httptest.NewServer(router)
			defer svr.Close()

			c := NewMetricHandler(svr.URL)
			_ = c.sendGaugeMetrics(tt.args.c)
		})
	}
}

func TestMetricSender_sendRand(t *testing.T) {
	type args struct {
		c collector.MetricCollectDTO
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Positive_test",
			args: args{
				c: collector.MetricCollectDTO{
					Rand: 1222.222,
				},
			},
			want: "/update/gauge/RandomValue/1222.222",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := http.NewServeMux()
			router.HandleFunc("/update/{type}/{name}/{value}", func(writer http.ResponseWriter, request *http.Request) {
				require.Equal(t, http.MethodPost, request.Method)
				assert.Equal(t, tt.want, request.RequestURI)
			})

			svr := httptest.NewServer(router)
			defer svr.Close()

			c := NewMetricHandler(svr.URL)
			_ = c.sendRand(tt.args.c)
		})
	}
}

func TestNewMetricHandler(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want *Sender
	}{
		{
			name: "Positive_test",
			args: args{
				url: "http://localhost",
			},
			want: &Sender{
				url: "http://localhost",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMetricHandler(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetricHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
