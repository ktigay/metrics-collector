package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ktigay/metrics-collector/internal/client/service/mocks"
	"github.com/ktigay/metrics-collector/internal/metric"
	"go.uber.org/zap"
)

func TestStatSenderService_SendStat(t *testing.T) {
	const (
		timeout = 100 * time.Millisecond
	)

	tests := []struct {
		name string
	}{
		{
			name: "Positive_test_SendStat",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)

			sender := mocks.NewMockStatSender(mockCtrl)
			sender.
				EXPECT().
				SendMetrics(gomock.All(), gomock.Any()).
				Do(func(_ []metric.Metrics, errChan chan<- error) {
					close(errChan)
				}).
				Times(1)

			handler := mocks.NewMockMetricsHandler(mockCtrl)
			handler.EXPECT().Processing(gomock.Any()).Times(1).Return([]metric.Metrics{})

			s := &StatSenderService{
				sender:   sender,
				handler:  handler,
				interval: 50 * time.Millisecond,
				logger:   zap.NewNop().Sugar(),
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			ch := make(chan []metric.Metrics, 1)
			defer close(ch)

			ch <- []metric.Metrics{}
			s.SendStat(ctx, ch)
		})
	}
}
