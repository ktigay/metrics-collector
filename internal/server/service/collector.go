package service

import (
	"fmt"
	"strconv"

	"github.com/ktigay/metrics-collector/internal/metric"
	e "github.com/ktigay/metrics-collector/internal/server/errors"
	"github.com/ktigay/metrics-collector/internal/server/storage"
)

// MetricStorage - интерфейс хранилища.
type MetricStorage interface {
	Save(m storage.Entity) error
	Find(t, n string) (*storage.Entity, error)
	Remove(t, n string) error
	All() []storage.Entity
	Backup() error
}

// MetricCollector - сборщик статистики.
type MetricCollector struct {
	storage MetricStorage
}

// NewMetricCollector - конструктор.
func NewMetricCollector(storage MetricStorage) *MetricCollector {
	return &MetricCollector{storage}
}

// Save - собирает статистику.
func (c *MetricCollector) Save(t, n string, v any) error {
	var (
		tp      metric.Type
		memItem *storage.Entity
		err     error
	)

	if tp, err = metric.ResolveType(t); err != nil {
		return e.ErrWrongType
	}

	if memItem, err = c.storage.Find(fmt.Sprint(tp), n); err != nil {
		return err
	}

	if memItem == nil {
		memItem = &storage.Entity{
			Key:  metric.Key(fmt.Sprint(tp), n),
			Name: n,
			Type: tp,
		}
	}

	switch tp {
	case metric.TypeCounter:
		var val int64
		switch vt := v.(type) {
		case string:
			if val, err = strconv.ParseInt(vt, 10, 64); err != nil {
				return e.ErrWrongValue
			}
		case int64:
			val = vt
		default:
			return e.ErrInvalidValueType
		}
		memItem.Delta = memItem.Delta + val
	case metric.TypeGauge:
		switch vt := v.(type) {
		case string:
			if memItem.Value, err = strconv.ParseFloat(vt, 64); err != nil {
				return e.ErrWrongValue
			}
		case float64:
			memItem.Value = vt
		default:
			return e.ErrInvalidValueType
		}
	}

	return c.storage.Save(*memItem)
}

// All - возвращает все записи.
func (c *MetricCollector) All() []storage.Entity {
	return c.storage.All()
}

// Find - находит запись по ключу.
func (c *MetricCollector) Find(t, n string) (*storage.Entity, error) {
	var (
		entity *storage.Entity
		err    error
	)
	if _, err = metric.ResolveType(t); err != nil {
		return nil, e.ErrWrongType
	}

	if entity, err = c.storage.Find(t, n); err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, e.ErrValueNotFound
	}

	return entity, nil
}

// Remove удаление записи по ключу.
func (c *MetricCollector) Remove(t, n string) error {
	if _, err := metric.ResolveType(t); err != nil {
		return e.ErrWrongType
	}
	return c.storage.Remove(t, n)
}

// Backup бэкап данных.
func (c *MetricCollector) Backup() error {
	return c.storage.Backup()
}
