package client

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ktigay/metrics-collector/internal/client/collector"
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricSender_sendCounter(t *testing.T) {
	type args struct {
		c collector.MetricCollectDTO
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Positive_test",
			args: args{
				c: collector.MetricCollectDTO{
					Counter: 15,
				},
			},
			want:    "/update/",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				require.Equal(t, http.MethodPost, request.Method)
				assert.Equal(t, tt.want, request.RequestURI)

				writer.WriteHeader(http.StatusCreated)
			}))
			defer svr.Close()

			c := NewSender(svr.URL, false)
			err := c.sendCounter(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("sendCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMetricSender_sendGaugeMetrics(t *testing.T) {
	type args struct {
		c collector.MetricCollectDTO
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Positive_test_check_request",
			args: args{
				c: collector.MetricCollectDTO{
					MemStats: map[metric.GaugeMetric]float64{
						metric.Alloc: 12.345,
					},
				},
			},
			want:    "/update/",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				require.Equal(t, http.MethodPost, request.Method)
				assert.Equal(t, tt.want, request.RequestURI)

				writer.WriteHeader(http.StatusCreated)
			}))
			defer svr.Close()

			c := NewSender(svr.URL, false)
			err := c.sendGaugeMetrics(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("sendGaugeMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetricSender_sendRand(t *testing.T) {
	type args struct {
		c collector.MetricCollectDTO
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Positive_test",
			args: args{
				c: collector.MetricCollectDTO{
					Rand: 1222.222,
				},
			},
			want:    "/update/",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				require.Equal(t, http.MethodPost, request.Method)
				assert.Equal(t, tt.want, request.RequestURI)

				writer.WriteHeader(http.StatusCreated)
			}))
			defer svr.Close()

			c := NewSender(svr.URL, false)
			err := c.sendRand(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("sendRand() error = %v, wantErr %v", err, tt.wantErr)
			}
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
			if got := NewSender(tt.args.url, false); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSender() = %v, want %v", got, tt.want)
			}
		})
	}
}
