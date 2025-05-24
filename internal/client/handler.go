// Package client агент.
package client

import (
	"io"
	"net/http"

	"github.com/ktigay/metrics-collector/internal/client/collector"
	"github.com/ktigay/metrics-collector/internal/compress"
	"github.com/ktigay/metrics-collector/internal/log"
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/retry"
	"go.uber.org/zap"
)

const (
	contentType  = "application/json"
	compressType = compress.Gzip
)

// Sender хендлер.
type Sender struct {
	url          string
	batchEnabled bool
}

// NewSender конструктор.
func NewSender(url string, batchEnabled bool) *Sender {
	return &Sender{
		url:          url,
		batchEnabled: batchEnabled,
	}
}

// SendMetrics отправляет метрики на сервер.
func (mh *Sender) SendMetrics(c collector.MetricCollectDTO) {
	if mh.batchEnabled {
		handler := func(policy retry.RetPolicy) error {
			log.AppLogger.Debugf("sending metrics to collector retries %d, err %v", policy.Retries(), policy.LastError())
			return mh.sendBatch(c)
		}
		if err := retry.Ret(handler); err != nil {
			log.AppLogger.Info("client.SendBatchMetrics error", zap.Error(err))
		}
		return
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
			log.AppLogger.Info("client.SendMetrics error", zap.Error(err))
		}
	}
}

func (mh *Sender) sendGaugeMetrics(c collector.MetricCollectDTO) (err error) {
	for n, m := range c.MemStats {
		mm := makeMetrics(metric.TypeGauge, string(n), m)
		if _, err = mh.post(mh.url+"/update/", mm); err != nil {
			break
		}
	}
	return err
}

func (mh *Sender) sendRand(c collector.MetricCollectDTO) error {
	mm := makeMetrics(metric.TypeGauge, metric.RandomValue, c.Rand)
	_, err := mh.post(mh.url+"/update/", mm)
	return err
}

func (mh *Sender) sendCounter(c collector.MetricCollectDTO) error {
	mm := makeMetrics(metric.TypeCounter, metric.PollCount, c.Counter)
	_, err := mh.post(mh.url+"/update/", mm)
	return err
}

func (mh *Sender) sendBatch(c collector.MetricCollectDTO) error {
	metrics := make([]metric.Metrics, 0, len(c.MemStats)+2)
	for n, m := range c.MemStats {
		mm := makeMetrics(metric.TypeGauge, string(n), m)
		metrics = append(metrics, mm)
	}
	metrics = append(metrics, makeMetrics(metric.TypeGauge, metric.RandomValue, c.Rand))
	metrics = append(metrics, makeMetrics(metric.TypeCounter, metric.PollCount, c.Counter))

	_, err := mh.post(mh.url+"/updates/", metrics)

	return err
}

func (mh *Sender) post(url string, body any) ([]byte, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
	)

	if req, err = compress.NewJSONRequest(
		http.MethodPost,
		url,
		compressType,
		body,
	); err != nil {
		return nil, err
	}

	if resp, err = compress.NewClient().Do(req); err != nil {
		return nil, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.AppLogger.Error("client.post error", zap.Error(err))
		}
	}()

	return io.ReadAll(resp.Body)
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
