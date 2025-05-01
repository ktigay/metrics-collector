package client

import (
	"bytes"
	"encoding/json"
	"github.com/ktigay/metrics-collector/internal/client/collector"
	"github.com/ktigay/metrics-collector/internal/metric"
	"log"
	"net/http"
)

const (
	contentType = "application/json"
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
	mh.sendGaugeMetrics(c)
	mh.sendRand(c)
	mh.sendCounter(c)
}

func (mh *Sender) sendGaugeMetrics(c collector.MetricCollectDTO) {
	for n, m := range c.MemStats {
		mh.post(mh.url+"/update/", metric.TypeGauge, string(n), m)
	}
}

func (mh *Sender) sendRand(c collector.MetricCollectDTO) {
	mh.post(mh.url+"/update/", metric.TypeGauge, metric.RandomValue, c.Rand)
}

func (mh *Sender) sendCounter(c collector.MetricCollectDTO) {
	mh.post(mh.url+"/update/", metric.TypeCounter, metric.PollCount, c.Counter)
}

func (mh *Sender) post(url string, t metric.Type, id string, v any) {

	m := makeMetrics(t, id, v)
	b, err := json.Marshal(m)

	if err != nil {
		log.Println(err)
		return
	}

	resp, err := http.Post(url, contentType, bytes.NewReader(b))

	if err != nil {
		log.Print(err)
	}
	if resp != nil && resp.StatusCode != http.StatusOK {
		log.Printf("Status code is not OK %d", resp.StatusCode)
	}

	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()
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
