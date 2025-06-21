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
		timeout = 1 * time.Second
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

			sender.EXPECT().SendMetrics(gomock.All(), gomock.Any()).Times(1)

			s := &StatSenderService{
				sender: sender,
				logger: zap.NewNop().Sugar(),
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			ch := make(chan []metric.Metrics)
			defer close(ch)

			go s.SendStat(ctx, ch)
			ch <- []metric.Metrics{}
		})
	}
}
