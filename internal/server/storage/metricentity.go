package storage

import (
	"strconv"

	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/server/errors"
)

// MetricEntity - сущность для сохранения в storage.
type MetricEntity struct {
	Key   string      `json:"key"`
	Type  metric.Type `json:"type"`
	Name  string      `json:"name"`
	Delta int64       `json:"delta"`
	Value float64     `json:"value"`
}

// ValueByType возвращает значение в зависимости от типа.
func (e *MetricEntity) ValueByType() any {
	switch e.Type {
	case metric.TypeCounter:
		return e.Delta
	case metric.TypeGauge:
		return e.Value
	}
	return nil
}

// ToMetrics мап сущности в дто.
func (e *MetricEntity) ToMetrics() metric.Metrics {
	m := metric.Metrics{
		ID:    e.Name,
		MType: string(e.Type),
	}

	switch e.Type {
	case metric.TypeCounter:
		m.Delta = &e.Delta
	case metric.TypeGauge:
		m.Value = &e.Value
	}

	return m
}

// AppendValue прибавляет/присваивает значение в зависимости от типа метрики.
func (e *MetricEntity) AppendValue(v any) error {
	var err error
	switch e.Type {
	case metric.TypeCounter:
		var val int64
		switch vt := v.(type) {
		case string:
			if val, err = strconv.ParseInt(vt, 10, 64); err != nil {
				return errors.ErrWrongValue
			}
		case int64:
			val = vt
		default:
			return errors.ErrInvalidValueType
		}
		e.Delta += val
	case metric.TypeGauge:
		switch vt := v.(type) {
		case string:
			if e.Value, err = strconv.ParseFloat(vt, 64); err != nil {
				return errors.ErrWrongValue
			}
		case float64:
			e.Value = vt
		default:
			return errors.ErrInvalidValueType
		}
	}

	return nil
}
