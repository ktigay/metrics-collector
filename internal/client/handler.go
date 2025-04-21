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
func (ms *MetricHandler) SendMetrics(c collector.MetricCollectDTO) {
	ms.sendGaugeMetrics(c)
	ms.sendRand(c)
	ms.sendCounter(c)
}

func (ms *MetricHandler) sendGaugeMetrics(c collector.MetricCollectDTO) {
	for n, m := range c.MemStats {
		func() {
			ms.post(fmt.Sprintf(ms.url+"/update/%s/%s/%v", metric.TypeGauge, n, m))
		}()
	}
}

func (ms *MetricHandler) sendRand(c collector.MetricCollectDTO) {
	ms.post(fmt.Sprintf(ms.url+"/update/%s/%s/%v", metric.TypeGauge, metric.RandomValue, c.Rand))
}

func (ms *MetricHandler) sendCounter(c collector.MetricCollectDTO) {
	ms.post(fmt.Sprintf(ms.url+"/update/%s/%s/%d", metric.TypeCounter, metric.PollCount, c.Counter))
}

func (ms *MetricHandler) post(url string) {
	body := strings.NewReader("")
	resp, err := http.Post(url, contentType, body)
	if err != nil {
		log.Print(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
}
