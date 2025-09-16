// Package grpc сервер.
package grpc

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ktigay/metrics-collector/internal/contracts"
	"github.com/ktigay/metrics-collector/internal/metric"
	e "github.com/ktigay/metrics-collector/internal/server/errors"
)

// Collector Интерфейс сборщика статистики.
//
//go:generate mockgen -destination=./mocks/mock_collector.go -package=mocks github.com/ktigay/metrics-collector/internal/server/handler/grpc Collector
type Collector interface {
	Find(ctx context.Context, t, n string) (*metric.Metrics, error)
	Save(ctx context.Context, mt metric.Metrics) (*metric.Metrics, error)
	SaveAll(ctx context.Context, mt []metric.Metrics) error
}

// MetricGrpcHandler структура с обработчиками запросов.
type MetricGrpcHandler struct {
	contracts.UnimplementedMetricsServiceServer
	collector Collector
	logger    *zap.SugaredLogger
}

// NewMetricGrpcHandler конструктор.
func NewMetricGrpcHandler(collector Collector, logger *zap.SugaredLogger) *MetricGrpcHandler {
	return &MetricGrpcHandler{
		collector: collector,
		logger:    logger,
	}
}

// GetMetrics возвращает метрику.
func (mh *MetricGrpcHandler) GetMetrics(ctx context.Context, req *contracts.GetMetricsRequest) (*contracts.GetMetricsResponse, error) {
	var (
		resp    contracts.GetMetricsResponse
		metrics *metric.Metrics
		err     error
	)

	if metrics, err = mh.collector.Find(ctx, req.GetType(), req.GetId()); err != nil {
		if errors.Is(err, e.ErrValueNotFound) {
			return nil, status.Error(codes.NotFound, "metric not found")
		}

		return nil, status.Errorf(codes.Internal, "failed to find metrics: %v", err)
	}

	resp.Metrics = metric.MapFromMetrics(metrics)

	return &resp, nil
}

// UpdateMetrics сохранение метрики.
func (mh *MetricGrpcHandler) UpdateMetrics(ctx context.Context, req *contracts.UpdateMetricsRequest) (*emptypb.Empty, error) {
	metrics := metric.MapToMetrics(req.Metrics)

	if _, err := mh.collector.Save(ctx, metrics); err != nil {
		mh.logger.Errorln("failed to update metrics: ", zap.Error(err))

		return nil, status.Errorf(codes.Internal, "failed to save metrics: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// BatchUpdateMetrics сохранение группы метрик.
func (mh *MetricGrpcHandler) BatchUpdateMetrics(ctx context.Context, req *contracts.BatchUpdateMetricsRequest) (*emptypb.Empty, error) {
	metrics := make([]metric.Metrics, 0, len(req.Metrics))
	for _, m := range req.Metrics {
		metrics = append(metrics, metric.MapToMetrics(m))
	}

	if err := mh.collector.SaveAll(ctx, metrics); err != nil {
		mh.logger.Errorln("failed to update metrics: ", zap.Error(err))

		return nil, status.Errorf(codes.Internal, "failed to save metrics: %v", err)
	}

	return &emptypb.Empty{}, nil
}
