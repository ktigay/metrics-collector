package main

import (
	"github.com/ktigay/metrics-collector/internal/client"
	"github.com/ktigay/metrics-collector/internal/client/collector"
	"sync"
	"time"
)

func main() {
	if err := parseFlags(); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	c := collector.NewRuntimeMetricCollector()
	h := client.NewMetricHandler(config.ServerProtocol + "://" + config.ServerHost)
	go pollStat(c)
	go sendStat(c, h)

	wg.Wait()
}

func pollStat(c *collector.RuntimeMetricCollector) {
	for {
		c.PollStat()
		time.Sleep(time.Duration(config.PollInterval) * time.Second)
	}
}

func sendStat(c *collector.RuntimeMetricCollector, h *client.MetricHandler) {
	for {
		time.Sleep(time.Duration(config.ReportInterval) * time.Second)
		h.SendMetrics(c.GetStat())
	}
}
