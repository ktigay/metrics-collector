package main

import (
	"github.com/ktigay/metrics-collector/internal/client"
	"github.com/ktigay/metrics-collector/internal/client/collector"
	"os"
	"time"
)

func main() {
	config, err := parseFlags(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	cl := collector.NewRuntimeMetricCollector()
	h := client.NewMetricHandler(config.ServerProtocol + "://" + config.ServerHost)
	stop := make(chan bool)
	defer func() {
		stop <- true
	}()
	go pollStat(stop, config, cl)

	sendStat(config, cl, h)
}

func pollStat(stop <-chan bool, config *Config, cl *collector.RuntimeMetricCollector) {
	for {
		select {
		default:
			cl.PollStat()
			time.Sleep(time.Duration(config.PollInterval) * time.Second)
		case <-stop:
			return
		}
	}
}

func sendStat(config *Config, cl *collector.RuntimeMetricCollector, h *client.Sender) {
	for {
		time.Sleep(time.Duration(config.ReportInterval) * time.Second)
		h.SendMetrics(cl.GetStat())
	}
}
