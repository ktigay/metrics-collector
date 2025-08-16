package sender

import (
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	"github.com/ktigay/metrics-collector/internal/client/sender/mocks"
	"github.com/ktigay/metrics-collector/internal/metric"
)

func TestMetricSender_sendGaugeMetrics(t *testing.T) {
	type args struct {
		c []metric.Metrics
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Positive_test_check_request",
			args: args{
				c: []metric.Metrics{
					{
						Type: "gauge",
						ID:   "Alloc",
						Value: func() *float64 {
							x := 12.0
							return &x
						}(),
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

			c := NewMetricSender(transport, false, 1, zap.NewNop().Sugar())

			resultCh := make(chan error)
			c.send(tt.args.c, resultCh)

			for err := range resultCh {
				if (err != nil) != tt.wantErr {
					t.Errorf("sendGaugeMetrics() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestSender_sendBatch(t *testing.T) {
	type args struct {
		c []metric.Metrics
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
				c: []metric.Metrics{
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
			wantLen: 3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			transport := mocks.NewMockTransport(mockCtrl)
			var b []byte
			transport.EXPECT().SendBatch(gomock.Len(tt.wantLen)).Return(b, nil).Times(1)

			c := NewMetricSender(transport, tt.batchEnabled, 1, zap.NewNop().Sugar())

			resultCh := make(chan error)
			c.sendBatch(tt.args.c, resultCh)

			for err := range resultCh {
				if (err != nil) != tt.wantErr {
					t.Errorf("SendBatch() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
