package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/caarlos0/env/v6"
)

const (
	defaultServerHost     = "localhost:8080"
	defaultLogLevel       = "info"
	defaultReportInterval = 10
	defaultPollInterval   = 2
	defaultServerProtocol = "http"
)

// Config - конфигурация клиента.
type Config struct {
	ServerProtocol string
	ServerHost     string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	LogLevel       string `env:"LOG_LEVEL"`
}

func parseFlags(args []string) (*Config, error) {
	config := &Config{
		ServerProtocol: defaultServerProtocol,
	}

	flags := flag.NewFlagSet("agent flags", flag.ContinueOnError)

	flags.StringVar(&config.ServerHost, "a", defaultServerHost, "address and port to run server")
	flags.StringVar(&config.LogLevel, "l", defaultLogLevel, "log level")
	flags.IntVar(&config.ReportInterval, "r", defaultReportInterval, "interval between reports")
	flags.IntVar(&config.PollInterval, "p", defaultPollInterval, "interval between polls")

	err := flags.Parse(args)
	if err != nil {
		return nil, err
	}

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	config.ServerHost = strings.TrimSpace(config.ServerHost)

	if config.ServerHost == "" {
		return nil, fmt.Errorf("host flag is required")
	}
	if config.ReportInterval < 1 {
		return nil, fmt.Errorf("report interval flag is required")
	}
	if config.PollInterval < 1 {
		return nil, fmt.Errorf("poll interval flag is required")
	}

	return config, nil
}
