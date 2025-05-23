package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ktigay/metrics-collector/internal/client/collector"
	"github.com/ktigay/metrics-collector/internal/compress"
	"github.com/ktigay/metrics-collector/internal/log"
	"github.com/ktigay/metrics-collector/internal/metric"
	"go.uber.org/zap"
)

const (
	contentType  = "application/json"
	compressType = compress.Gzip
)

// Sender - хендлер.
type Sender struct {
	url string
}

// NewSender - конструктор.
func NewSender(url string) *Sender {
	return &Sender{
		url: url,
	}
}

// SendMetrics - отправляет метрики на сервер.
func (mh *Sender) SendMetrics(c collector.MetricCollectDTO) {
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

func (mh *Sender) sendGaugeMetrics(c collector.MetricCollectDTO) error {
	for n, m := range c.MemStats {
		_, err := mh.post(mh.url+"/update/", metric.TypeGauge, string(n), m)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mh *Sender) sendRand(c collector.MetricCollectDTO) error {
	_, err := mh.post(mh.url+"/update/", metric.TypeGauge, metric.RandomValue, c.Rand)
	return err
}

func (mh *Sender) sendCounter(c collector.MetricCollectDTO) error {
	_, err := mh.post(mh.url+"/update/", metric.TypeCounter, metric.PollCount, c.Counter)
	return err
}

func (mh *Sender) post(url string, t metric.Type, id string, v any) ([]byte, error) {
	var (
		m    metric.Metrics
		b    []byte
		err  error
		buff bytes.Buffer
		cw   *compress.Writer
		req  *http.Request
		resp *http.Response
	)

	m = makeMetrics(t, id, v)

	if b, err = json.Marshal(m); err != nil {
		return nil, err
	}
	if cw, err = compress.NewWriter(compressType, &buff); err != nil {
		return nil, err
	}
	if _, err = cw.Write(b); err != nil {
		return nil, err
	}
	if err = cw.Close(); err != nil {
		return nil, err
	}

	if req, err = http.NewRequest(http.MethodPost, url, &buff); err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Encoding", string(compressType))

	client := http.Client{}
	if resp, err = client.Do(req); err != nil {
		return nil, err
	}
	defer func() {
		if resp == nil {
			return
		}
		if err = resp.Body.Close(); err != nil {
			log.AppLogger.Error("client.post error", zap.Error(err))
		}
	}()

	if resp != nil && (resp.StatusCode > 300 || resp.StatusCode < 200) {
		return nil, fmt.Errorf("status code is not OK %d", resp.StatusCode)
	}

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
