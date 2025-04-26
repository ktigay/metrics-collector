package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

const (
	defaultServerHost = ":8080"
)

// Config - конфигурация сервера.
type Config struct {
	ServerHost string `env:"ADDRESS"`
}

func parseFlags(args []string) (*Config, error) {
	config := &Config{}
	flags := flag.NewFlagSet("server flags", flag.ContinueOnError)
	flags.StringVar(&config.ServerHost, "a", defaultServerHost, "address and port to run server")

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
