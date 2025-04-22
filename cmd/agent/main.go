package main

import (
	"github.com/ktigay/metrics-collector/internal/client"
	"github.com/ktigay/metrics-collector/internal/client/collector"
	"time"
)

func main() {
	if err := parseFlags(); err != nil {
		panic(err)
	}

	c := collector.NewRuntimeMetricCollector()
	h := client.NewMetricHandler(config.ServerProtocol + "://" + config.ServerHost)
	go pollStat(c)

	sendStat(c, h)
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
