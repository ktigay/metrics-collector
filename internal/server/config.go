package server

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

const (
	defaultServerHost      = ":8080"
	defaultLogLevel        = "info"
	defaultStoreInterval   = 300
	defaultFileStoragePath = "/tmp/metrics-db.json"
	defaultRestoreFlag     = true
	defaultDatabaseDSN     = ""
)

// Config - конфигурация сервера.
type Config struct {
	ServerHost      string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

// InitializeConfig инициализирует конфигурацию.
func InitializeConfig(args []string) (*Config, error) {
	var err error

	config := Config{}

	flags := flag.NewFlagSet("server flags", flag.ContinueOnError)

	flags.StringVar(&config.ServerHost, "a", defaultServerHost, "address and port to run server")
	flags.StringVar(&config.LogLevel, "l", defaultLogLevel, "log level")
	flags.IntVar(&config.StoreInterval, "i", defaultStoreInterval, "storage interval in seconds")
	flags.StringVar(&config.FileStoragePath, "f", defaultFileStoragePath, "file storage path")
	flags.BoolVar(&config.Restore, "r", defaultRestoreFlag, "restore data from storage")
	flags.StringVar(&config.DatabaseDSN, "d", defaultDatabaseDSN, "database DSN")

	if err = flags.Parse(args); err != nil {
		return nil, err
	}
	if err = env.Parse(&config); err != nil {
		return nil, err
	}
	if config.ServerHost == "" {
		return nil, fmt.Errorf("host flag is required")
	}

	return &config, nil
}
