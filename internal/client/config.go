package client

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	defaultServerHost     = "localhost:8080"
	defaultLogLevel       = "debug"
	defaultReportInterval = 10
	defaultPollInterval   = 2
	defaultServerProtocol = "http"
	defaultBatchEnabled   = false
	defaultHashKey        = ""
	defaultRateLimit      = 1
	defaultCryptoKey      = ""
)

// ConfigInterval интервал в секундах.
type ConfigInterval int64

// UnmarshalJSON десериализация из JSON.
func (i *ConfigInterval) UnmarshalJSON(bytes []byte) error {
	v, err := time.ParseDuration(strings.Trim(string(bytes), `"`))
	if err != nil {
		return err
	}
	*i = ConfigInterval(v.Seconds())
	return nil
}

// Config конфигурация клиента.
type Config struct {
	ServerProtocol string
	ServerHost     string         `env:"ADDRESS" json:"address"`
	LogLevel       string         `env:"LOG_LEVEL"`
	HashKey        string         `env:"KEY"`
	CryptoKey      string         `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigFile     string         `env:"CONFIG"`
	IPAddr         string         `json:"ip_addr"`
	BatchEnabled   bool           `env:"BATCH_ENABLED"`
	ReportInterval ConfigInterval `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval   ConfigInterval `env:"POLL_INTERVAL" json:"poll_interval"`
	RateLimit      int            `env:"RATE_LIMIT"`
}

// InitializeConfig инициализирует конфиг клиента.
func InitializeConfig(args []string) (*Config, error) {
	config := Config{}
	config.setDefaults()

	var (
		configFilePath string
		exists         bool
		err            error
	)

	if configFilePath, exists = os.LookupEnv("CONFIG"); !exists {
		for i, f := range args {
			if strings.HasPrefix(f, "-c=") {
				configFilePath = strings.TrimPrefix(f, "-c=")
				break
			} else if f == "-c" && len(args) > i+1 {
				configFilePath = args[i+1]
				break
			}
		}
	}

	if configFilePath != "" {
		if err = config.parseFromFile(configFilePath); err != nil {
			return nil, err
		}
	}

	if err = config.parseFromArgs(args); err != nil {
		return nil, err
	}

	if err = config.parseFromEnv(); err != nil {
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

func (c *Config) parseFromFile(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, c)
}

func (c *Config) parseFromArgs(args []string) error {
	flags := flag.NewFlagSet("agent flags", flag.ContinueOnError)

	flags.StringVar(&c.ServerHost, "a", c.ServerHost, "address and port to run server")
	flags.StringVar(&c.LogLevel, "lvl", c.LogLevel, "log level")
	flags.BoolVar(&c.BatchEnabled, "b", c.BatchEnabled, "enable batchEnabled request")
	flags.StringVar(&c.HashKey, "k", c.HashKey, "SHA256 hash key")
	flags.IntVar(&c.RateLimit, "l", c.RateLimit, "requests rate limit")
	flags.StringVar(&c.CryptoKey, "crypto-key", c.CryptoKey, "Public key path")

	flags.StringVar(&c.ConfigFile, "c", c.ConfigFile, "JSON config file path")

	var reportInterval, pollInterval int64
	flags.Int64Var(&reportInterval, "r", int64(c.ReportInterval), "interval between reports")

	flags.Int64Var(&pollInterval, "p", int64(c.PollInterval), "interval between polls")

	if err := flags.Parse(args); err != nil {
		return err
	}

	c.ReportInterval = ConfigInterval(reportInterval)
	c.PollInterval = ConfigInterval(pollInterval)

	return nil
}

func (c *Config) parseFromEnv() error {
	return env.Parse(c)
}

func (c *Config) setDefaults() {
	c.ServerProtocol = defaultServerProtocol
	c.ServerHost = defaultServerHost
	c.LogLevel = defaultLogLevel
	c.BatchEnabled = defaultBatchEnabled
	c.HashKey = defaultHashKey
	c.RateLimit = defaultRateLimit
	c.CryptoKey = defaultCryptoKey
	c.ReportInterval = defaultReportInterval
	c.PollInterval = defaultPollInterval
}
