package config

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
	defaultGRPCHost       = ":3000"
	defaultTransport      = "grpc"
)

// Interval интервал в секундах.
type Interval int64

// UnmarshalJSON десериализация из JSON.
func (i *Interval) UnmarshalJSON(bytes []byte) error {
	v, err := time.ParseDuration(strings.Trim(string(bytes), `"`))
	if err != nil {
		return err
	}
	*i = Interval(v.Seconds())
	return nil
}

// TransportProtocol тип протокола отправки сообщений.
type TransportProtocol string

// Типы протоколов.
var (
	TransportHTTP TransportProtocol = "http"
	TransportGRPC TransportProtocol = "grpc"
)

// Config конфигурация клиента.
type Config struct {
	ServerProtocol string
	ServerHost     string `env:"ADDRESS" json:"address"`
	LogLevel       string `env:"LOG_LEVEL"`
	HashKey        string `env:"KEY"`
	CryptoKey      string `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigFile     string `env:"CONFIG"`
	IPAddr         string `json:"ip_addr"`
	ServerGRPCHost string `env:"GRPC_ADDRESS" json:"grpc_address"`
	Transport      TransportProtocol
	BatchEnabled   bool     `env:"BATCH_ENABLED"`
	ReportInterval Interval `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval   Interval `env:"POLL_INTERVAL" json:"poll_interval"`
	RateLimit      int      `env:"RATE_LIMIT"`
}

// NewConfig конструктор.
func NewConfig(arguments []string) (*Config, error) {
	config := &Config{}

	handler := DefaultHandler{
		next: &FileHandler{
			arguments: arguments,
			next: &ArgumentsHandler{
				arguments: arguments,
				next: &EnvHandler{
					next: &ValidateHandler{},
				},
			},
		},
	}

	return handler.Handle(config)
}

// Handler интерфейс парсера конфигурации.
type Handler interface {
	Handle(*Config) (*Config, error)
}

// DefaultHandler дефолтные значения.
type DefaultHandler struct {
	next Handler
}

// Handle обработчик.
func (d *DefaultHandler) Handle(c *Config) (*Config, error) {
	c.ServerProtocol = defaultServerProtocol
	c.ServerHost = defaultServerHost
	c.LogLevel = defaultLogLevel
	c.BatchEnabled = defaultBatchEnabled
	c.HashKey = defaultHashKey
	c.RateLimit = defaultRateLimit
	c.CryptoKey = defaultCryptoKey
	c.ReportInterval = defaultReportInterval
	c.PollInterval = defaultPollInterval
	c.ServerGRPCHost = defaultGRPCHost
	c.Transport = defaultTransport

	return d.next.Handle(c)
}

// FileHandler конфиг из файла.
type FileHandler struct {
	next      Handler
	arguments []string
}

// Handle обработчик.
func (f *FileHandler) Handle(c *Config) (*Config, error) {
	var (
		err    error
		path   string
		exists bool
	)

	if path, exists = os.LookupEnv("CONFIG"); !exists {
		for i, arg := range f.arguments {
			if strings.HasPrefix(arg, "-c=") {
				path = strings.TrimPrefix(arg, "-c=")
				break
			} else if arg == "-c" && len(f.arguments) > i+1 {
				path = f.arguments[i+1]
				break
			}
		}
	}

	if path == "" {
		return f.next.Handle(c)
	}

	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	var content []byte
	if content, err = os.ReadFile(path); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(content, c); err != nil {
		return nil, err
	}

	return f.next.Handle(c)
}

// ArgumentsHandler конфиг из аргументов.
type ArgumentsHandler struct {
	next      Handler
	arguments []string
}

// Handle обработчик.
func (a *ArgumentsHandler) Handle(c *Config) (*Config, error) {
	flags := flag.NewFlagSet("agent flags", flag.ContinueOnError)

	flags.StringVar(&c.ServerHost, "a", c.ServerHost, "address and port of server")
	flags.StringVar(&c.LogLevel, "lvl", c.LogLevel, "log level")
	flags.BoolVar(&c.BatchEnabled, "b", c.BatchEnabled, "enable batchEnabled request")
	flags.StringVar(&c.HashKey, "k", c.HashKey, "SHA256 hash key")
	flags.IntVar(&c.RateLimit, "l", c.RateLimit, "requests rate limit")
	flags.StringVar(&c.CryptoKey, "crypto-key", c.CryptoKey, "Public key path")
	flags.StringVar(&c.ServerGRPCHost, "g", c.ServerGRPCHost, "address and port of grpc server")

	flags.StringVar(&c.ConfigFile, "c", c.ConfigFile, "JSON config file path")

	var reportInterval, pollInterval int64
	flags.Int64Var(&reportInterval, "r", int64(c.ReportInterval), "interval between reports")

	flags.Int64Var(&pollInterval, "p", int64(c.PollInterval), "interval between polls")

	if err := flags.Parse(a.arguments); err != nil {
		return nil, err
	}

	c.ReportInterval = Interval(reportInterval)
	c.PollInterval = Interval(pollInterval)

	return a.next.Handle(c)
}

// EnvHandler конфиг из переменных среды.
type EnvHandler struct {
	next Handler
}

// Handle обработчик.
func (e *EnvHandler) Handle(c *Config) (*Config, error) {
	if err := env.Parse(c); err != nil {
		return nil, err
	}

	return e.next.Handle(c)
}

// ValidateHandler валидация.
type ValidateHandler struct{}

// Handle обработчик.
func (e *ValidateHandler) Handle(c *Config) (*Config, error) {
	c.ServerHost = strings.TrimSpace(c.ServerHost)

	if c.ServerHost == "" {
		return nil, fmt.Errorf("host flag is required")
	}
	if c.ReportInterval < 1 {
		return nil, fmt.Errorf("report interval flag is required")
	}
	if c.PollInterval < 1 {
		return nil, fmt.Errorf("poll interval flag is required")
	}

	return c, nil
}
