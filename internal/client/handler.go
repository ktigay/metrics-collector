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

// MetricHandler - хендлер.
type MetricHandler struct {
	url string
}

// NewMetricHandler - конструктор.
func NewMetricHandler(url string) *MetricHandler {
	return &MetricHandler{
		url: url,
	}
}

// SendMetrics - отправляет метрики на сервер.
func (mh *MetricHandler) SendMetrics(c collector.MetricCollectDTO) {
	mh.sendGaugeMetrics(c)
	mh.sendRand(c)
	mh.sendCounter(c)
}

func (mh *MetricHandler) sendGaugeMetrics(c collector.MetricCollectDTO) {
	for n, m := range c.MemStats {
		func() {
			mh.post(fmt.Sprintf(mh.url+"/update/%s/%s/%v", metric.TypeGauge, n, m))
		}()
	}
}

func (mh *MetricHandler) sendRand(c collector.MetricCollectDTO) {
	mh.post(fmt.Sprintf(mh.url+"/update/%s/%s/%v", metric.TypeGauge, metric.RandomValue, c.Rand))
}

func (mh *MetricHandler) sendCounter(c collector.MetricCollectDTO) {
	mh.post(fmt.Sprintf(mh.url+"/update/%s/%s/%d", metric.TypeCounter, metric.PollCount, c.Counter))
}

func (mh *MetricHandler) post(url string) {
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
