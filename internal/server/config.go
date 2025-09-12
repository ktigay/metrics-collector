package server

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

// Config конфигурация сервера.
type Config struct {
	ServerHost      string         `env:"ADDRESS" json:"address"`
	LogLevel        string         `env:"LOG_LEVEL" json:"log_level"`
	FileStoragePath string         `env:"FILE_STORAGE_PATH" json:"store_file"`
	DatabaseDSN     string         `env:"DATABASE_DSN" json:"database_dsn"`
	DatabaseDriver  string         `env:"DATABASE_DRIVER" json:"database_driver"`
	HashKey         string         `env:"KEY" json:"hash_key"`
	CryptoKey       string         `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigFile      string         `env:"CONFIG"`
	TrustedSubnet   string         `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	StoreInterval   ConfigInterval `env:"STORE_INTERVAL" json:"store_interval"`
	Restore         bool           `env:"RESTORE" json:"restore"`
}

// IsUseSQLDB использовать БД SQL.
func (c *Config) IsUseSQLDB() bool {
	return c.DatabaseDSN != "" && c.DatabaseDriver != ""
}

// InitializeConfig инициализирует конфигурацию.
func InitializeConfig(args []string) (*Config, error) {
	config := Config{}
	config.setDefaults()

	var (
		err            error
		configFilePath string
		exists         bool
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
		err = config.parseFromFile(configFilePath)
		if err != nil {
			return nil, err
		}
	}

	if err = config.parseFromArgs(args); err != nil {
		return nil, err
	}

	if err = config.parseFromEnv(); err != nil {
		return nil, err
	}

	if config.ServerHost == "" {
		return nil, fmt.Errorf("host flag is required")
	}

	return &config, nil
}

func (c *Config) parseFromArgs(args []string) error {
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

	if err := flags.Parse(args); err != nil {
		return err
	}

	c.StoreInterval = ConfigInterval(storeInterval)

	return nil
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

func (c *Config) parseFromEnv() error {
	return env.Parse(c)
}

func (c *Config) setDefaults() {
	c.DatabaseDriver = defaultDatabaseDriver
	c.ServerHost = defaultServerHost
	c.LogLevel = defaultLogLevel
	c.StoreInterval = defaultStoreInterval
	c.FileStoragePath = defaultFileStoragePath
	c.Restore = defaultRestoreFlag
	c.DatabaseDSN = defaultDatabaseDSN
	c.HashKey = defaultHashKey
	c.CryptoKey = defaultCryptoKey
}
