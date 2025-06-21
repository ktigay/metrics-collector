package collector

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ktigay/metrics-collector/internal/client/collector/mocks"
	"github.com/ktigay/metrics-collector/internal/metric"
	"go.uber.org/zap"
)

func TestIntervalPoller_PollStat(t *testing.T) {
	type args struct {
		timeout time.Duration
	}
	type fields struct {
		interval time.Duration
	}
	tests := []struct {
		name      string
		args      args
		fields    fields
		wantTimes int
	}{
		{
			name: "Positive_test_PollStat_MultipleCalls",
			args: args{
				timeout: 200 * time.Millisecond,
			},
			fields: fields{
				interval: 60 * time.Millisecond,
			},
			wantTimes: 3,
		},
		{
			name: "PollStat_Abort_With_Timeout",
			args: args{
				timeout: 50 * time.Millisecond,
			},
			fields: fields{
				interval: 100 * time.Millisecond,
			},
			wantTimes: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			sg := mocks.NewMockStatGetter(mockCtrl)
			sg.EXPECT().GetStat().Times(tt.wantTimes)

			m := &IntervalPoller{
				source:   sg,
				interval: tt.fields.interval,
				logger:   zap.NewNop().Sugar(),
			}

			ch := make(chan []metric.Metrics)
			defer close(ch)

			ctx, cancel := context.WithTimeout(context.Background(), tt.args.timeout)
			defer cancel()

			go func() {
				for {
					<-ch
				}
			}()
			m.PollStat(ctx, ch)
		})
	}
}
