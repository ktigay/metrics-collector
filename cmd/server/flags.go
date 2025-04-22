package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

var config = struct {
	ServerHost string `env:"ADDRESS"`
}{}

func parseFlags() error {
	flag.StringVar(&config.ServerHost, "a", ":8080", "address and port to run server")

	flag.Parse()

	if err := env.Parse(&config); err != nil {
		return err
	}

	if config.ServerHost == "" {
		return fmt.Errorf("host flag is required")
	}

	return nil
}
