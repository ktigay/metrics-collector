package service

import (
	"fmt"
	"github.com/ktigay/metrics-collector/internal/log"
	"strconv"

	"github.com/ktigay/metrics-collector/internal/metric"
	e "github.com/ktigay/metrics-collector/internal/server/errors"
	"github.com/ktigay/metrics-collector/internal/server/storage"
)

// MetricStorage - интерфейс хранилища.
type MetricStorage interface {
	Save(m storage.MetricEntity) error
	Find(t, n string) (*storage.MetricEntity, error)
	Remove(t, n string) error
	All() ([]storage.MetricEntity, error)
}

type BackupStorage interface {
	Backup() error
	Restore() error
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
		memItem *storage.MetricEntity
		err     error
	)

	if tp, err = metric.ResolveType(t); err != nil {
		return e.ErrWrongType
	}

	if memItem, err = c.storage.Find(fmt.Sprint(tp), n); err != nil {
		return err
	}

	if memItem == nil {
		memItem = &storage.MetricEntity{
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
func (c *MetricCollector) All() ([]storage.MetricEntity, error) {
	return c.storage.All()
}

// Find - находит запись по ключу.
func (c *MetricCollector) Find(t, n string) (*storage.MetricEntity, error) {
	var (
		entity *storage.MetricEntity
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
	switch t := c.storage.(type) {
	case BackupStorage:
		return t.Backup()
	}

	log.AppLogger.Debug("storage not supported backups")
	return nil
}

// Restore восстановление данных.
func (c *MetricCollector) Restore() error {
	switch t := c.storage.(type) {
	case BackupStorage:
		return t.Restore()
	}

	log.AppLogger.Debug("storage not supported restores")
	return nil
}
