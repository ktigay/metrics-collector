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

func TestRuntimeMetricCollector_PollStat(t *testing.T) {
	const (
		timeout = 1 * time.Second
	)

	tests := []struct {
		name string
	}{
		{
			name: "Positive_test_PushStat",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			storage := mocks.NewMockStorageInterface(mockCtrl)

			storage.EXPECT().Save(gomock.All()).Times(1)

			c := StatSaverService{
				storage: storage,
				logger:  zap.NewNop().Sugar(),
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			ch := make(chan []metric.Metrics)
			defer close(ch)

			go c.PushStat(ctx, ch)
			ch <- []metric.Metrics{}
		})
	}
}
