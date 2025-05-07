package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

const (
	defaultServerHost      = ":8080"
	defaultLogLevel        = "info"
	defaultStoreInterval   = 300
	defaultFileStoragePath = "./cache/storage.txt"
	defaultRestoreFlag     = true
)

// Config - конфигурация сервера.
type Config struct {
	ServerHost      string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

func parseFlags(args []string) (*Config, error) {
	config := &Config{}
	flags := flag.NewFlagSet("server flags", flag.ContinueOnError)

	flags.StringVar(&config.ServerHost, "a", defaultServerHost, "address and port to run server")
	flags.StringVar(&config.LogLevel, "l", defaultLogLevel, "log level")
	flags.IntVar(&config.StoreInterval, "i", defaultStoreInterval, "storage interval in seconds")
	flags.StringVar(&config.FileStoragePath, "f", defaultFileStoragePath, "file storage path")
	flags.BoolVar(&config.Restore, "r", defaultRestoreFlag, "restore data from storage")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	if config.ServerHost == "" {
		return nil, fmt.Errorf("host flag is required")
	}

	return config, nil
}
