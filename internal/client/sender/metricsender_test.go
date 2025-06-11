package sender

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ktigay/metrics-collector/internal/client/collector"
	"github.com/ktigay/metrics-collector/internal/client/sender/mocks"
	"github.com/ktigay/metrics-collector/internal/metric"
)

func TestMetricSender_sendCounter(t *testing.T) {
	type args struct {
		c collector.MetricCollectDTO
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Positive_test",
			args: args{
				c: collector.MetricCollectDTO{
					Counter: 15,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			transport := mocks.NewMockTransport(mockCtrl)
			var b []byte
			transport.EXPECT().Send(gomock.Any()).Return(b, nil).Times(1)

			c := NewMetricSender(transport, false)
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
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			transport := mocks.NewMockTransport(mockCtrl)
			var b []byte
			transport.EXPECT().Send(gomock.Any()).Return(b, nil).Times(1)

			c := NewMetricSender(transport, false)
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
		wantErr bool
	}{
		{
			name: "Positive_test",
			args: args{
				c: collector.MetricCollectDTO{
					Rand: 1222.222,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			transport := mocks.NewMockTransport(mockCtrl)
			var b []byte
			transport.EXPECT().Send(gomock.Any()).Return(b, nil).Times(1)

			c := NewMetricSender(transport, false)
			err := c.sendRand(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("sendRand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSender_sendBatch(t *testing.T) {
	type args struct {
		c collector.MetricCollectDTO
	}
	tests := []struct {
		name         string
		batchEnabled bool
		args         args
		wantLen      int
		wantErr      bool
	}{
		{
			name:         "Positive_test",
			batchEnabled: true,
			args: args{
				c: collector.MetricCollectDTO{
					MemStats: map[metric.GaugeMetric]float64{
						metric.Alloc: 12.345,
						metric.GCSys: 22.345,
					},
					Counter: 15,
					Rand:    1222.222,
				},
			},
			wantLen: 4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			transport := mocks.NewMockTransport(mockCtrl)
			var b []byte
			transport.EXPECT().SendBatch(gomock.Len(tt.wantLen)).Return(b, nil).Times(1)

			c := NewMetricSender(transport, tt.batchEnabled)
			err := c.sendBatch(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
