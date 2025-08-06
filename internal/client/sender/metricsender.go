// Package client агент.
package sender

import (
	"github.com/ktigay/metrics-collector/internal/metric"
	"go.uber.org/zap"
)

// Transport транспорт.
//
//go:generate mockgen -destination=./mocks/mock_transport.go -package=mocks github.com/ktigay/metrics-collector/internal/client/sender Transport
type Transport interface {
	Send(body metric.Metrics) ([]byte, error)
	SendBatch(body []metric.Metrics) ([]byte, error)
}

// MetricSender хендлер.
type MetricSender struct {
	transport    Transport
	batchEnabled bool
	rateLimit    int
	logger       *zap.SugaredLogger
}

// SendMetrics отправляет метрики на сервер.
func (mh *MetricSender) SendMetrics(metrics []metric.Metrics, resultCh chan<- error) {
	if mh.batchEnabled {
		mh.sendBatch(metrics, resultCh)
		return
	}

	mh.send(metrics, resultCh)
}

func (mh *MetricSender) send(metrics []metric.Metrics, resultCh chan<- error) {
	rateLimit := mh.rateLimit
	if rateLimit <= 0 {
		rateLimit = 1
	}

	jobs := make(chan metric.Metrics, len(metrics))

	errCh := make(chan error, rateLimit)

	for i := 1; i <= rateLimit; i++ {
		go mh.worker(i, jobs, errCh)
	}

	for _, m := range metrics {
		jobs <- m
	}
	close(jobs)

	go func() {
		defer close(errCh)
		defer close(resultCh)

		for i := 0; i < rateLimit; i++ {
			resultCh <- <-errCh
		}
	}()
}

func (mh *MetricSender) sendBatch(metrics []metric.Metrics, resultCh chan<- error) {
	go func() {
		defer close(resultCh)
		_, err := mh.transport.SendBatch(metrics)
		resultCh <- err
	}()
}

func (mh *MetricSender) worker(thread int, jobs <-chan metric.Metrics, errChan chan<- error) {
	for j := range jobs {
		mh.logger.Debugf("Sending metrics, thread #%d", thread)
		if _, err := mh.transport.Send(j); err != nil {
			errChan <- err
			return
		}
	}
	errChan <- nil
}

// NewMetricSender конструктор.
func NewMetricSender(transport Transport, batchEnabled bool, rateLimit int, logger *zap.SugaredLogger) *MetricSender {
	return &MetricSender{
		transport:    transport,
		batchEnabled: batchEnabled,
		rateLimit:    rateLimit,
		logger:       logger,
	}
}
