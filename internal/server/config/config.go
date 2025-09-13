package config

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	defaultServerHost      = ":8080"
	defaultLogLevel        = "debug"
	defaultStoreInterval   = 300
	defaultFileStoragePath = "/tmp/metrics-db.json"
	defaultRestoreFlag     = false
	defaultDatabaseDSN     = ""
	defaultDatabaseDriver  = "pgx"
	defaultHashKey         = ""
	defaultCryptoKey       = ""
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

// Config конфигурация сервера.
type Config struct {
	ServerHost      string   `env:"ADDRESS" json:"address"`
	LogLevel        string   `env:"LOG_LEVEL" json:"log_level"`
	FileStoragePath string   `env:"FILE_STORAGE_PATH" json:"store_file"`
	DatabaseDSN     string   `env:"DATABASE_DSN" json:"database_dsn"`
	DatabaseDriver  string   `env:"DATABASE_DRIVER" json:"database_driver"`
	HashKey         string   `env:"KEY" json:"hash_key"`
	CryptoKey       string   `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigFile      string   `env:"CONFIG"`
	TrustedSubnet   string   `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	StoreInterval   Interval `env:"STORE_INTERVAL" json:"store_interval"`
	Restore         bool     `env:"RESTORE" json:"restore"`
}

// IsUseSQLDB использовать БД SQL.
func (c *Config) IsUseSQLDB() bool {
	return c.DatabaseDSN != "" && c.DatabaseDriver != ""
}

// NewConfig конструктор.
func NewConfig(arguments []string) (*Config, error) {
	config := &Config{}

	handler := DefaultHandler{
		next: &FileHandler{
			arguments: arguments,
			next: &ArgumentsHandler{
				arguments: arguments,
				next:      &EnvHandler{},
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
	c.DatabaseDriver = defaultDatabaseDriver
	c.ServerHost = defaultServerHost
	c.LogLevel = defaultLogLevel
	c.StoreInterval = defaultStoreInterval
	c.FileStoragePath = defaultFileStoragePath
	c.Restore = defaultRestoreFlag
	c.DatabaseDSN = defaultDatabaseDSN
	c.HashKey = defaultHashKey
	c.CryptoKey = defaultCryptoKey

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
	flags := flag.NewFlagSet("server flags", flag.ContinueOnError)

	flags.StringVar(&c.ServerHost, "a", c.ServerHost, "address and port to run server")
	flags.StringVar(&c.LogLevel, "lvl", c.LogLevel, "log level")
	flags.StringVar(&c.FileStoragePath, "f", c.FileStoragePath, "file storage path")
	flags.BoolVar(&c.Restore, "r", c.Restore, "restore data from storage")
	flags.StringVar(&c.DatabaseDSN, "d", c.DatabaseDSN, "database DSN")
	flags.StringVar(&c.HashKey, "k", c.HashKey, "SHA256 hash key")
	flags.StringVar(&c.CryptoKey, "crypto-key", c.CryptoKey, "Private key path")
	flags.StringVar(&c.TrustedSubnet, "t", c.TrustedSubnet, "Trusted subnet")

	flags.StringVar(&c.ConfigFile, "c", c.ConfigFile, "JSON config file path")

	var storeInterval int64
	flags.Int64Var(&storeInterval, "i", int64(c.StoreInterval), "storage interval in seconds")

	if err := flags.Parse(a.arguments); err != nil {
		return nil, err
	}

	c.StoreInterval = Interval(storeInterval)

	return a.next.Handle(c)
}

// EnvHandler конфиг из переменных среды.
type EnvHandler struct{}

// Handle обработчик.
func (e *EnvHandler) Handle(c *Config) (*Config, error) {
	if err := env.Parse(c); err != nil {
		return nil, err
	}
	return c, nil
}
