package service

import (
	"fmt"

	"github.com/ktigay/metrics-collector/internal/log"
	"github.com/ktigay/metrics-collector/internal/metric"
	e "github.com/ktigay/metrics-collector/internal/server/errors"
	"github.com/ktigay/metrics-collector/internal/server/storage"
)

// MetricStorage интерфейс хранилища.
type MetricStorage interface {
	Upsert(m storage.MetricEntity) error
	Find(t, n string) (*storage.MetricEntity, error)
	Remove(t, n string) error
	All() ([]storage.MetricEntity, error)
}

// BackupStorage интерфейс для работы со снапшотами.
type BackupStorage interface {
	Backup() error
	Restore() error
}

// BatchMetricStorage интерфейс для работы с батчами.
type BatchMetricStorage interface {
	UpsertAll(mt []storage.MetricEntity) error
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

	memItem = &storage.MetricEntity{
		Key:  metric.Key(fmt.Sprint(tp), n),
		Name: n,
		Type: tp,
	}
	if err = memItem.AppendValue(v); err != nil {
		return err
	}

	return c.storage.Upsert(*memItem)
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

// SaveAll сохраняет батч.
func (c *MetricCollector) SaveAll(mt []metric.Metrics) error {
	var err error
	entities := make([]storage.MetricEntity, 0, len(mt))

	for _, m := range mt {
		var t metric.Type
		if t, err = metric.ResolveType(m.MType); err != nil {
			return e.ErrWrongType
		}

		en := storage.MetricEntity{
			Key:  metric.Key(m.MType, m.ID),
			Name: m.ID,
			Type: t,
		}
		if err = en.AppendValue(m.ValueByType()); err != nil {
			return err
		}
		entities = append(entities, en)
	}

	switch t := c.storage.(type) {
	case BatchMetricStorage:
		return t.UpsertAll(entities)
	default:
		for _, en := range entities {
			if err = c.storage.Upsert(en); err != nil {
				return err
			}
		}
		return nil
	}
}
