package main

import (
	"github.com/gorilla/mux"
	"log"

	"github.com/ktigay/metrics-collector/internal/server"
	"github.com/ktigay/metrics-collector/internal/server/collector"
	"github.com/ktigay/metrics-collector/internal/server/storage"

	"net/http"
)

func main() {

	if err := parseFlags(); err != nil {
		log.Fatal(err)
	}

	c := collector.NewMetricCollector(storage.NewMemStorage())
	s := server.NewServer(c)

	router := mux.NewRouter()

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			next.ServeHTTP(w, r)
		})
	})

	router.HandleFunc("/update/{type}/{name}/{value}", s.CollectHandler).Methods(http.MethodPost)
	router.HandleFunc("/value/{type}/{name}", s.GetValueHandler).Methods(http.MethodGet)
	router.HandleFunc("/", s.GetAllHandler).Methods(http.MethodGet)

	err := http.ListenAndServe(config.ServerHost, router)
	if err != nil {
		panic(err)
	}
}
