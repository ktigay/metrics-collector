package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ktigay/metrics-collector/internal/client/collector"
	cio "github.com/ktigay/metrics-collector/internal/client/io"
	"github.com/ktigay/metrics-collector/internal/compress"
	"github.com/ktigay/metrics-collector/internal/metric"
	"io"
	"log"
	"net/http"
)

const (
	contentType  = "application/json"
	compressType = compress.Gzip
)

// Sender - хендлер.
type Sender struct {
	url string
}

// NewMetricHandler - конструктор.
func NewMetricHandler(url string) *Sender {
	return &Sender{
		url: url,
	}
}

// SendMetrics - отправляет метрики на сервер.
func (mh *Sender) SendMetrics(c collector.MetricCollectDTO) {
	if err := mh.sendGaugeMetrics(c); err != nil {
		log.Println(err)
		return
	}
	if err := mh.sendRand(c); err != nil {
		log.Println(err)
		return
	}
	if err := mh.sendCounter(c); err != nil {
		log.Println(err)
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

	m := makeMetrics(t, id, v)
	b, err := json.Marshal(m)

	if err != nil {
		return nil, err
	}

	var bb bytes.Buffer
	cw, err := cio.NewCompressWriter(compressType, &bb)
	if err != nil {
		return nil, err
	}
	if _, err = cw.Write(b); err != nil {
		return nil, err
	}
	if err = cw.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, &bb)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Encoding", string(compressType))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp != nil && (resp.StatusCode > 300 || resp.StatusCode < 200) {
		return nil, fmt.Errorf("status code is not OK %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func makeMetrics(t metric.Type, id string, v any) metric.Metrics {
	var delta int64
	var val float64

	switch t := v.(type) {
	case int64:
		delta = t
	case float64:
		val = t
	}

	return metric.Metrics{
		ID:    id,
		MType: string(t),
		Delta: &delta,
		Value: &val,
	}
}
