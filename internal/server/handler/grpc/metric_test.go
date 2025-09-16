package grpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ktigay/metrics-collector/internal/contracts"
	"github.com/ktigay/metrics-collector/internal/metric"
	e "github.com/ktigay/metrics-collector/internal/server/errors"
	"github.com/ktigay/metrics-collector/internal/server/handler/grpc/mocks"
)

func TestMetricGrpcHandler_BatchUpdateMetrics(t *testing.T) {
	type args struct {
		req *contracts.BatchUpdateMetricsRequest
	}
	tests := []struct {
		respErr     error
		args        args
		name        string
		wantErrCode codes.Code
		wantErr     bool
	}{
		{
			name: "Success",
			args: args{
				req: &contracts.BatchUpdateMetricsRequest{
					Metrics: []*contracts.Metrics{},
				},
			},
			wantErr:     false,
			wantErrCode: codes.OK,
		},
		{
			name: "Internal_Error",
			args: args{
				req: &contracts.BatchUpdateMetricsRequest{
					Metrics: []*contracts.Metrics{},
				},
			},
			respErr:     fmt.Errorf("some error"),
			wantErr:     true,
			wantErrCode: codes.Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			collector := mocks.NewMockCollector(mockCtrl)
			collector.EXPECT().SaveAll(gomock.Any(), gomock.Any()).Return(tt.respErr).Times(1)

			mh := &MetricGrpcHandler{
				collector: collector,
				logger:    zap.NewNop().Sugar(),
			}
			_, err := mh.BatchUpdateMetrics(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchUpdateMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			statusErr, _ := status.FromError(err)
			if statusErr.Code() != tt.wantErrCode {
				t.Errorf("BatchUpdateMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMetricGrpcHandler_GetMetrics(t *testing.T) {
	type args struct {
		req *contracts.GetMetricsRequest
	}
	tests := []struct {
		respErr     error
		args        args
		name        string
		wantErrCode codes.Code
		wantErr     bool
	}{
		{
			name: "Success",
			args: args{
				req: &contracts.GetMetricsRequest{},
			},
			wantErr:     false,
			wantErrCode: codes.OK,
		},
		{
			name: "Internal_Error",
			args: args{
				req: &contracts.GetMetricsRequest{},
			},
			respErr:     fmt.Errorf("some error"),
			wantErr:     true,
			wantErrCode: codes.Internal,
		},
		{
			name: "NotFound_Error",
			args: args{
				req: &contracts.GetMetricsRequest{},
			},
			respErr:     e.ErrValueNotFound,
			wantErr:     true,
			wantErrCode: codes.NotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			collector := mocks.NewMockCollector(mockCtrl)
			collector.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(&metric.Metrics{}, tt.respErr).Times(1)

			mh := &MetricGrpcHandler{
				collector: collector,
				logger:    zap.NewNop().Sugar(),
			}
			_, err := mh.GetMetrics(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			statusErr, _ := status.FromError(err)
			if statusErr.Code() != tt.wantErrCode {
				t.Errorf("GetMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMetricGrpcHandler_UpdateMetrics(t *testing.T) {
	type args struct {
		req *contracts.UpdateMetricsRequest
	}
	tests := []struct {
		respErr     error
		args        args
		name        string
		wantErrCode codes.Code
		wantErr     bool
	}{
		{
			name: "Success",
			args: args{
				req: &contracts.UpdateMetricsRequest{
					Metrics: &contracts.Metrics{},
				},
			},
			wantErr:     false,
			wantErrCode: codes.OK,
		},
		{
			name: "Internal_Error",
			args: args{
				req: &contracts.UpdateMetricsRequest{
					Metrics: &contracts.Metrics{},
				},
			},
			respErr:     fmt.Errorf("some error"),
			wantErr:     true,
			wantErrCode: codes.Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			collector := mocks.NewMockCollector(mockCtrl)
			collector.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil, tt.respErr).Times(1)

			mh := &MetricGrpcHandler{
				collector: collector,
				logger:    zap.NewNop().Sugar(),
			}
			_, err := mh.UpdateMetrics(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			statusErr, _ := status.FromError(err)
			if statusErr.Code() != tt.wantErrCode {
				t.Errorf("UpdateMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
