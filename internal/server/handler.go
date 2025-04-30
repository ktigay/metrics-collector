package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/server/storage"
	"go.uber.org/zap"
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
	logger    *zap.Logger
}

// NewServer - конструктор.
func NewServer(collector CollectorInterface, logger *zap.Logger) *Server {
	return &Server{collector, logger}
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
	if _, err := fmt.Fprintf(w, "%v", entity.GetValue()); err != nil {
		c.logger.Sugar().Errorln("Failed to write response", zap.Error(err))
	}
}

// GetAllHandler - обработчик для получения списка метрик.
func (c *Server) GetAllHandler(w http.ResponseWriter, _ *http.Request) {
	metrics := c.collector.GetAll()

	names := make([]string, 0, len(metrics))
	for _, m := range metrics {
		names = append(names, m.Name)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(names); err != nil {
		c.logger.Sugar().Errorln("Failed to write response", zap.Error(err))
	}
}
