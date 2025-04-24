package client

import (
	"fmt"
	"github.com/ktigay/metrics-collector/internal/client/collector"
	"github.com/ktigay/metrics-collector/internal/metric"
	"log"
	"net/http"
	"strings"
)

const (
	contentType = "text/plain"
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
		mh.post(fmt.Sprintf(mh.url+"/update/%s/%s/%v", metric.TypeGauge, n, m))
	}
}

func (mh *Sender) sendRand(c collector.MetricCollectDTO) {
	mh.post(fmt.Sprintf(mh.url+"/update/%s/%s/%v", metric.TypeGauge, metric.RandomValue, c.Rand))
}

func (mh *Sender) sendCounter(c collector.MetricCollectDTO) {
	mh.post(fmt.Sprintf(mh.url+"/update/%s/%s/%d", metric.TypeCounter, metric.PollCount, c.Counter))
}

func (mh *Sender) post(url string) {
	body := strings.NewReader("")
	resp, err := http.Post(url, contentType, body)
	if err != nil {
		log.Print(err)
	}
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()
}
