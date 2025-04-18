package server

import (
	"errors"
	"github.com/ktigay/metrics-collector/internal/server/collector"
	"github.com/ktigay/metrics-collector/internal/server/metric/item"
	"net/http"
)

type Server struct {
	collector *collector.MetricCollector
}

func NewServer(collector *collector.MetricCollector) Server {
	return Server{collector}
}

func (c *Server) CollectHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "text/plain")

	if r.Method != http.MethodPost {
		return
	}

	m, err := item.ParseFromPath(r.URL.Path)

	if err != nil {
		switch {
		case errors.Is(err, item.ErrorInvalidName):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	_ = c.collector.Save(m)

	w.WriteHeader(http.StatusOK)
}
