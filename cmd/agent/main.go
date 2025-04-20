package main

import (
	"github.com/ktigay/metrics-collector/internal/client"
	"github.com/ktigay/metrics-collector/internal/client/collector"
	"sync"
	"time"
)

const (
	serverURL      = "http://127.0.0.1:8080"
	pollInterval   = 2
	reportInterval = 10
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	c := collector.NewRuntimeMetricCollector()
	h := client.NewMetricHandler(serverURL)
	go pollStat(c)
	go sendStat(c, h)

	wg.Wait()
}

func pollStat(c *collector.RuntimeMetricCollector) {
	for {
		c.PollStat()
		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}

func sendStat(c *collector.RuntimeMetricCollector, h *client.MetricHandler) {
	for {
		time.Sleep(time.Duration(reportInterval) * time.Second)
		h.SendMetrics(c.GetStat())
	}
}
