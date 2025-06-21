// Package handler сервер.
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/server/errors"
	"github.com/ktigay/metrics-collector/internal/server/repository"
	"go.uber.org/zap"
)

var errStatusMap = map[error]int{
	errors.ErrWrongType:     http.StatusBadRequest,
	errors.ErrWrongValue:    http.StatusBadRequest,
	errors.ErrValueNotFound: http.StatusNotFound,
}

func statusFromError(err error) int {
	if st, ok := errStatusMap[err]; ok {
		return st
	}
	return http.StatusInternalServerError
}

// CollectorInterface Интерфейс сборщика статистики.
//
//go:generate mockgen -destination=./mocks/mock_collector.go -package=mocks github.com/ktigay/metrics-collector/internal/server/handler CollectorInterface
type CollectorInterface interface {
	Save(ctx context.Context, mt metric.Metrics) error
	All(ctx context.Context) ([]repository.MetricEntity, error)
	Find(ctx context.Context, t, n string) (*metric.Metrics, error)
	Remove(ctx context.Context, t, n string) error
	SaveAll(ctx context.Context, mt []metric.Metrics) error
}

// MetricHandler структура с обработчиками запросов.
type MetricHandler struct {
	collector CollectorInterface
	logger    *zap.SugaredLogger
}

// NewMetricHandler конструктор.
func NewMetricHandler(collector CollectorInterface, logger *zap.SugaredLogger) *MetricHandler {
	return &MetricHandler{
		collector: collector,
		logger:    logger,
	}
}

// CollectHandler обработчик для сборка метрик.
func (mh *MetricHandler) CollectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	mt := metric.Metrics{
		ID:   vars["name"],
		Type: vars["type"],
	}
	if err := mt.SetValueByType(vars["value"]); err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	if err := mh.collector.Save(r.Context(), mt); err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}
}

// GetValueHandler обработчик для получения значения метрики.
func (mh *MetricHandler) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var (
		err error
		m   *metric.Metrics
	)

	if m, err = mh.collector.Find(r.Context(), vars["type"], vars["name"]); err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	if _, err = fmt.Fprintf(w, "%v", m.ValueByType()); err != nil {
		mh.logger.Errorln("Failed to write response", zap.Error(err))
	}
}

// GetAllHandler обработчик для получения списка метрик.
func (mh *MetricHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	metrics, _ := mh.collector.All(r.Context())

	names := make([]string, 0, len(metrics))
	for _, m := range metrics {
		names = append(names, m.Name)
	}

	if err := json.NewEncoder(w).Encode(names); err != nil {
		mh.logger.Errorln("Failed to write response", zap.Error(err))
	}
}

// UpdateJSONHandler обработчик обновления метрики из json-строки.
func (mh *MetricHandler) UpdateJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")

	var (
		m   metric.Metrics
		err error
		mm  *metric.Metrics
	)

	if err = json.NewDecoder(r.Body).Decode(&m); err != nil {
		mh.logger.Errorln("Failed to write response", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	if err = mh.collector.Save(ctx, m); err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	if mm, err = mh.collector.Find(ctx, m.Type, m.ID); err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	if err = json.NewEncoder(w).Encode(mm); err != nil {
		mh.logger.Errorln("Failed to write response", zap.Error(err))
	}
}

// UpdatesJSONHandler обработчик обновления метрик из json-строки.
func (mh *MetricHandler) UpdatesJSONHandler(w http.ResponseWriter, r *http.Request) {
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
		mh.logger.Errorln("Failed to write response", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = mh.collector.SaveAll(r.Context(), m); err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}
}

// GetJSONValueHandler возвращает структуру в виде json-строки.
func (mh *MetricHandler) GetJSONValueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")

	var (
		m   metric.Metrics
		err error
		mm  *metric.Metrics
	)

	if err = json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mm, err = mh.collector.Find(r.Context(), m.Type, m.ID)
	if err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	if err = json.NewEncoder(w).Encode(mm); err != nil {
		mh.logger.Errorln("Failed to write response", zap.Error(err))
	}
}
