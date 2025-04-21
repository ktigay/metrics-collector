package main

import (
	"flag"
	"fmt"
)

var config = struct {
	ServerHost string
}{}

func parseFlags() error {
	flag.StringVar(&config.ServerHost, "a", ":8080", "address and port to run server")

	flag.Parse()

	if config.ServerHost == "" {
		return fmt.Errorf("host flag is required")
	}

	return nil
}
