package main

import (
	"github.com/ktigay/metrics-collector/internal/server"
	"github.com/ktigay/metrics-collector/internal/server/collector"
	"github.com/ktigay/metrics-collector/internal/server/storage"
	"net/http"
)

func main() {
	c := collector.NewMetricCollector(storage.NewMemStorage())
	server := server.NewServer(c)

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", server.CollectHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
