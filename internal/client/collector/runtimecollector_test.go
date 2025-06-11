package collector

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ktigay/metrics-collector/internal/client/collector/mocks"
)

func TestRuntimeMetricCollector_PollStat(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Positive_test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			storage := mocks.NewMockStorageInterface(mockCtrl)

			storage.EXPECT().Save(gomock.All()).Times(1)

			c := &RuntimeMetricCollector{
				storage: storage,
			}

			c.PollStat(context.TODO())
		})
	}
}
