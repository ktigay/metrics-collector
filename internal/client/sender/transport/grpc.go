package transport

import (
	"context"

	"go.uber.org/zap"

	"github.com/ktigay/metrics-collector/internal/contracts/mapper"

	"github.com/ktigay/metrics-collector/internal/contracts"
	"github.com/ktigay/metrics-collector/internal/metric"
)

// GRPCClient grpc клиент.
type GRPCClient struct {
	conn   contracts.MetricsServiceClient
	logger *zap.SugaredLogger
}

// NewGRPCClient конструктор.
func NewGRPCClient(conn contracts.MetricsServiceClient, logger *zap.SugaredLogger) *GRPCClient {
	return &GRPCClient{
		conn:   conn,
		logger: logger,
	}
}

// Send отправка одной метрики.
func (g *GRPCClient) Send(body metric.Metrics) ([]byte, error) {
	req := contracts.UpdateMetricsRequest{
		Metrics: mapper.MapFromMetrics(&body),
	}

	resp, err := g.conn.UpdateMetrics(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return []byte(resp.String()), nil
}

// SendBatch отправка батча.
func (g *GRPCClient) SendBatch(body []metric.Metrics) ([]byte, error) {
	mt := make([]*contracts.Metrics, 0, len(body))
	for _, m := range body {
		mt = append(mt, mapper.MapFromMetrics(&m))
	}
	req := contracts.BatchUpdateMetricsRequest{
		Metrics: mt,
	}

	resp, err := g.conn.BatchUpdateMetrics(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return []byte(resp.String()), nil
}
