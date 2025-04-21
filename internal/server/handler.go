package server

import (
	"github.com/gorilla/mux"
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/server/collector"
	"github.com/ktigay/metrics-collector/internal/server/metric/item"
	"net/http"
	"strconv"
	"strings"
)

// Server - структура с обработчиками запросов.
type Server struct {
	collector *collector.MetricCollector
}

// NewServer - конструктор.
func NewServer(collector *collector.MetricCollector) Server {
	return Server{collector}
}

// CollectHandler - обработчик для сборка метрик.
func (c *Server) CollectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := metric.ResolveType(vars["type"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := strconv.ParseFloat(vars["value"], 64); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m := item.MetricDTO{
		Type:  t,
		Name:  vars["name"],
		Value: vars["value"],
	}

	if err := c.collector.Save(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetValueHandler - обработчик для получения значения метрики.
func (c *Server) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if _, err := metric.ResolveType(vars["type"]); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	dto, err := c.collector.FindByKey(metric.GetKey(vars["type"], vars["name"]))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if dto == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(dto.GetValue()))
}

func (c *Server) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	metrics := c.collector.GetAll()

	names := make([]string, 0)
	for _, m := range metrics {
		names = append(names, m.Name)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strings.Join(names, "\n")))
}
