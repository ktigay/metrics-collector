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
func (mh *MetricSender) SendMetrics(metrics []metric.Metrics) error {
	if mh.batchEnabled {
		return mh.sendBatch(metrics)
	}

	return mh.send(metrics)
}

func (mh *MetricSender) send(metrics []metric.Metrics) (err error) {
	rateLimit := mh.rateLimit
	if rateLimit <= 0 {
		rateLimit = 1
	}

	jobs := make(chan metric.Metrics, len(metrics))

	errChan := make(chan error, rateLimit)
	defer close(errChan)

	for i := 1; i <= rateLimit; i++ {
		go mh.worker(i, jobs, errChan)
	}

	for _, m := range metrics {
		jobs <- m
	}
	close(jobs)

	for i := 0; i < rateLimit; i++ {
		if e := <-errChan; e != nil {
			err = e
		}
	}

	return err
}

func (mh *MetricSender) sendBatch(metrics []metric.Metrics) error {
	_, err := mh.transport.SendBatch(metrics)

	return err
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
