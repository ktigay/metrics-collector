// Package server сервер.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ktigay/metrics-collector/internal/log"
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/server/db"
	"github.com/ktigay/metrics-collector/internal/server/errors"
	"github.com/ktigay/metrics-collector/internal/server/storage"
	"go.uber.org/zap"
)

var errStatusMap = map[error]int{
	errors.ErrWrongType:     http.StatusBadRequest,
	errors.ErrWrongValue:    http.StatusBadRequest,
	errors.ErrValueNotFound: http.StatusNotFound,
}

const (
	pingTimeout = 2 * time.Second
)

func statusFromError(err error) int {
	if st, ok := errStatusMap[err]; ok {
		return st
	}
	return http.StatusInternalServerError
}

// CollectorInterface Интерфейс сборщика статистики.
type CollectorInterface interface {
	Save(ctx context.Context, t, n string, v any) error
	All(ctx context.Context) ([]storage.MetricEntity, error)
	Find(ctx context.Context, t, n string) (*storage.MetricEntity, error)
	Remove(ctx context.Context, t, n string) error
	SaveAll(ctx context.Context, mt []metric.Metrics) error
}

// Server структура с обработчиками запросов.
type Server struct {
	collector CollectorInterface
}

// NewServer конструктор.
func NewServer(collector CollectorInterface) *Server {
	return &Server{collector}
}

// CollectHandler обработчик для сборка метрик.
func (c *Server) CollectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := c.collector.Save(r.Context(), vars["type"], vars["name"], vars["value"]); err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetValueHandler обработчик для получения значения метрики.
func (c *Server) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var (
		err    error
		entity *storage.MetricEntity
	)

	if entity, err = c.collector.Find(r.Context(), vars["type"], vars["name"]); err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err = fmt.Fprintf(w, "%v", entity.ValueByType()); err != nil {
		log.AppLogger.Errorln("Failed to write response", zap.Error(err))
	}
}

// GetAllHandler обработчик для получения списка метрик.
func (c *Server) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	metrics, _ := c.collector.All(r.Context())

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
		entity *storage.MetricEntity
	)

	if err = json.NewDecoder(r.Body).Decode(&m); err != nil {
		log.AppLogger.Errorln("Failed to write response", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	if err = c.collector.Save(ctx, m.MType, m.ID, m.ValueByType()); err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	if entity, err = c.collector.Find(ctx, m.MType, m.ID); err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	w.WriteHeader(http.StatusOK)

	um := entity.ToMetrics()
	if err = json.NewEncoder(w).Encode(um); err != nil {
		log.AppLogger.Errorln("Failed to write response", zap.Error(err))
	}
}

// UpdatesJSONHandler обработчик обновления метрик из json-строки.
func (c *Server) UpdatesJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")

	var (
		m   []metric.Metrics
		err error
	)

	if err = json.NewDecoder(r.Body).Decode(&m); err != nil {
		log.AppLogger.Errorln("Failed to write response", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = c.collector.SaveAll(r.Context(), m); err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	w.WriteHeader(http.StatusOK)
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
		entity *storage.MetricEntity
	)

	if err = json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	entity, err = c.collector.Find(r.Context(), m.MType, m.ID)
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

// Ping пинг соединения с БД.
func (c *Server) Ping(w http.ResponseWriter, r *http.Request) {
	if db.MasterDB == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), pingTimeout)
	defer cancel()

	if err = db.MasterDB.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.AppLogger.Errorf("Failed to connect to MasterDB %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
