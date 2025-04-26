package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/server/storage"
	"net/http"
	"strconv"
)

// CollectorInterface - Интерфейс сборщика статистики.
type CollectorInterface interface {
	Save(t metric.Type, n string, v any) error
	GetAll() map[string]*storage.Entity
	FindByKey(key string) (*storage.Entity, error)
}

// Server - структура с обработчиками запросов.
type Server struct {
	collector CollectorInterface
}

// NewServer - конструктор.
func NewServer(collector CollectorInterface) *Server {
	return &Server{collector}
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

	if err := c.collector.Save(t, vars["name"], vars["value"]); err != nil {
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

	entity, err := c.collector.FindByKey(metric.GetKey(vars["type"], vars["name"]))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if entity == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%v", entity.GetValue())))
}

// GetAllHandler - обработчик для получения списка метрик.
func (c *Server) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	metrics := c.collector.GetAll()

	names := make([]string, len(metrics))
	for _, m := range metrics {
		names = append(names, m.Name)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(names)
}
