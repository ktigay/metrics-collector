package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
)

var config = struct {
	ServerProtocol string
	ServerHost     string
	ReportInterval int
	PollInterval   int
}{
	ServerProtocol: "http",
}

var flags = []string{"a", "r", "p"}

func parseFlags() error {
	flag.StringVar(&config.ServerHost, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&config.ReportInterval, "r", 10, "interval between reports")
	flag.IntVar(&config.PollInterval, "p", 2, "interval between polls")

	flag.Parse()

	if config.ServerHost == "" {
		return fmt.Errorf("host flag is required")
	}
	if config.ReportInterval < 1 {
		return fmt.Errorf("report interval flag is required")
	}
	if config.PollInterval < 1 {
		return fmt.Errorf("poll interval flag is required")
	}

	for _, v := range os.Args[1:] {
		if !slices.Contains(flags, v[1:2]) {
			return fmt.Errorf("unknown flag: %s", v)
		}
	}
	return nil
}
