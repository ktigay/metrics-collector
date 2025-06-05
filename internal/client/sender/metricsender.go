// Package client агент.
package sender

import (
	"github.com/ktigay/metrics-collector/internal/client/collector"
	"github.com/ktigay/metrics-collector/internal/metric"
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
}

// SendMetrics отправляет метрики на сервер.
func (mh *MetricSender) SendMetrics(c collector.MetricCollectDTO) error {
	var lastError error

	if mh.batchEnabled {
		return mh.sendBatch(c)
	}

	errChan := make(chan error, 3)

	go func() {
		if err := mh.sendGaugeMetrics(c); err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	go func() {
		if err := mh.sendRand(c); err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	go func() {
		if err := mh.sendCounter(c); err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	// Wait for all operations to complete
	for i := 0; i < 3; i++ {
		if err := <-errChan; err != nil {
			lastError = err
		}
	}

	return lastError
}

func (mh *MetricSender) sendGaugeMetrics(c collector.MetricCollectDTO) (err error) {
	for n, m := range c.MemStats {
		mm := makeMetrics(metric.TypeGauge, string(n), m)
		if _, err = mh.transport.Send(mm); err != nil {
			break
		}
	}
	return err
}

func (mh *MetricSender) sendRand(c collector.MetricCollectDTO) error {
	mm := makeMetrics(metric.TypeGauge, metric.RandomValue, c.Rand)
	_, err := mh.transport.Send(mm)
	return err
}

func (mh *MetricSender) sendCounter(c collector.MetricCollectDTO) error {
	mm := makeMetrics(metric.TypeCounter, metric.PollCount, c.Counter)
	_, err := mh.transport.Send(mm)
	return err
}

func (mh *MetricSender) sendBatch(c collector.MetricCollectDTO) error {
	metrics := make([]metric.Metrics, 0, len(c.MemStats)+2)
	for n, m := range c.MemStats {
		mm := makeMetrics(metric.TypeGauge, string(n), m)
		metrics = append(metrics, mm)
	}
	metrics = append(metrics, makeMetrics(metric.TypeGauge, metric.RandomValue, c.Rand))
	metrics = append(metrics, makeMetrics(metric.TypeCounter, metric.PollCount, c.Counter))

	_, err := mh.transport.SendBatch(metrics)

	return err
}

// NewMetricSender конструктор.
func NewMetricSender(transport Transport, batchEnabled bool) *MetricSender {
	return &MetricSender{
		transport:    transport,
		batchEnabled: batchEnabled,
	}
}

func makeMetrics(t metric.Type, id string, v any) metric.Metrics {
	var (
		delta int64
		val   float64
	)

	switch mt := v.(type) {
	case int64:
		delta = mt
	case float64:
		val = mt
	}

	return metric.Metrics{
		ID:    id,
		MType: string(t),
		Delta: &delta,
		Value: &val,
	}
}
