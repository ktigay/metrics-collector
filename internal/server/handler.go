package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ktigay/metrics-collector/internal/log"
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/server/errors"
	"github.com/ktigay/metrics-collector/internal/server/storage"
	"go.uber.org/zap"
)

var errStatusMap = map[error]int{
	errors.ErrWrongType:     http.StatusBadRequest,
	errors.ErrTypeNotFound:  http.StatusBadRequest,
	errors.ErrWrongValue:    http.StatusBadRequest,
	errors.ErrValueNotFound: http.StatusNotFound,
}

func statusFromError(err error) int {
	if st, ok := errStatusMap[err]; ok {
		return st
	}
	return http.StatusInternalServerError
}

// CollectorInterface - Интерфейс сборщика статистики.
type CollectorInterface interface {
	Save(t, n string, v any) error
	All() []storage.Entity
	Find(t, n string) (*storage.Entity, error)
	Remove(t, n string) error
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

	if err := c.collector.Save(vars["type"], vars["name"], vars["value"]); err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetValueHandler - обработчик для получения значения метрики.
func (c *Server) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var (
		err    error
		entity *storage.Entity
	)

	if entity, err = c.collector.Find(vars["type"], vars["name"]); err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err = fmt.Fprintf(w, "%v", entity.ValueByType()); err != nil {
		log.AppLogger.Errorln("Failed to write response", zap.Error(err))
	}
}

// GetAllHandler - обработчик для получения списка метрик.
func (c *Server) GetAllHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	metrics := c.collector.All()

	names := make([]string, 0, len(metrics))
	for _, m := range metrics {
		names = append(names, m.Name)
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(names); err != nil {
		log.AppLogger.Errorln("Failed to write response", zap.Error(err))
	}
}

// UpdateJSONHandler обработчик обновления метрики из json-строки.
func (c *Server) UpdateJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")

	var (
		m      metric.Metrics
		err    error
		entity *storage.Entity
	)

	if err = json.NewDecoder(r.Body).Decode(&m); err != nil {
		log.AppLogger.Errorln("Failed to write response", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = c.collector.Save(m.MType, m.ID, m.ValueByType()); err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	if entity, err = c.collector.Find(m.MType, m.ID); err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	w.WriteHeader(http.StatusOK)

	um := entity.ToMetrics()
	if err = json.NewEncoder(w).Encode(um); err != nil {
		log.AppLogger.Errorln("Failed to write response", zap.Error(err))
	}
}

// GetJSONValueHandler возвращает структуру в виде json-строки.
func (c *Server) GetJSONValueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")

	var (
		m      metric.Metrics
		err    error
		entity *storage.Entity
	)

	if err = json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	entity, err = c.collector.Find(m.MType, m.ID)
	if err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	w.WriteHeader(http.StatusOK)

	um := entity.ToMetrics()
	if err = json.NewEncoder(w).Encode(um); err != nil {
		log.AppLogger.Errorln("Failed to write response", zap.Error(err))
	}
}
