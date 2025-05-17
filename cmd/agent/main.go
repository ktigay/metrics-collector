package main

import (
	"log"
	"os"
	"time"

	"github.com/ktigay/metrics-collector/internal/client"
	"github.com/ktigay/metrics-collector/internal/client/collector"
	ilog "github.com/ktigay/metrics-collector/internal/log"
)

func main() {
	config, err := parseFlags(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	if err = ilog.Initialize(config.LogLevel); err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func() {
		if err := ilog.AppLogger.Sync(); err != nil {
			log.Printf("can't sync logger: %v", err)
		}
	}()

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
