package server

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ktigay/metrics-collector/internal/log"
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/server/storage"
)

// CollectorInterface - Интерфейс сборщика статистики.
type CollectorInterface interface {
	Save(t metric.Type, n string, v any) error
	GetAll() []storage.Entity
	FindByKey(key string) (*storage.Entity, error)
	RemoveByKey(key string) error
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

	entity, err := c.collector.FindByKey(metric.Key(vars["type"], vars["name"]))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if entity == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprintf(w, "%v", entity.ValueByType()); err != nil {
		log.SugaredLogger.Errorln("Failed to write response", zap.Error(err))
	}
}

// GetAllHandler - обработчик для получения списка метрик.
func (c *Server) GetAllHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	metrics := c.collector.GetAll()

	names := make([]string, 0, len(metrics))
	for _, m := range metrics {
		names = append(names, m.Name)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(names); err != nil {
		log.SugaredLogger.Errorln("Failed to write response", zap.Error(err))
	}
}

// UpdateJSONHandler обработчик обновления метрики из json-строки.
func (c *Server) UpdateJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")

	var m metric.Metrics

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		log.SugaredLogger.Errorln("Failed to write response", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t, err := metric.ResolveType(m.MType)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.SugaredLogger.Errorln("resolve type error", zap.Error(err))
		return
	}

	if err := c.collector.Save(t, m.ID, m.ValueByType()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	entity, err := c.collector.FindByKey(m.Key())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if entity == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

	um := storage.MapEntityToMetrics(*entity)
	if err := json.NewEncoder(w).Encode(um); err != nil {
		log.SugaredLogger.Errorln("Failed to write response", zap.Error(err))
	}
}

// GetJSONValueHandler возвращает структуру в виде json-строки.
func (c *Server) GetJSONValueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")

	var m metric.Metrics
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := metric.ResolveType(m.MType); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.SugaredLogger.Errorln("resolve type error", zap.Error(err))
		return
	}

	entity, err := c.collector.FindByKey(m.Key())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if entity == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

	um := storage.MapEntityToMetrics(*entity)
	if err := json.NewEncoder(w).Encode(um); err != nil {
		log.SugaredLogger.Errorln("Failed to write response", zap.Error(err))
	}
}
