package client

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
	defaultBatchEnabled   = false
	defaultHashKey        = ""
	defaultRateLimit      = 1
)

// Config конфигурация клиента.
type Config struct {
	ServerProtocol string
	ServerHost     string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	LogLevel       string `env:"LOG_LEVEL"`
	BatchEnabled   bool   `env:"BATCH_ENABLED"`
	HashKey        string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

// InitializeConfig инициализирует конфиг клиента.
func InitializeConfig(args []string) (*Config, error) {
	config := Config{
		ServerProtocol: defaultServerProtocol,
	}

	flags := flag.NewFlagSet("agent flags", flag.ContinueOnError)

	flags.StringVar(&config.ServerHost, "a", defaultServerHost, "address and port to run server")
	flags.StringVar(&config.LogLevel, "lvl", defaultLogLevel, "log level")
	flags.IntVar(&config.ReportInterval, "r", defaultReportInterval, "interval between reports")
	flags.IntVar(&config.PollInterval, "p", defaultPollInterval, "interval between polls")
	flags.BoolVar(&config.BatchEnabled, "b", defaultBatchEnabled, "enable batchEnabled request")
	flags.StringVar(&config.HashKey, "k", defaultHashKey, "SHA256 hash key")
	flags.IntVar(&config.RateLimit, "l", defaultRateLimit, "requests rate limit")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	if err := env.Parse(&config); err != nil {
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

	return &config, nil
}
