package collector

import (
	"errors"
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/server/storage"
	"strconv"
)

// StorageInterface - интерфейс хранилища.
type StorageInterface interface {
	Save(m *storage.Entity) error
	GetAll() map[string]*storage.Entity
	FindByKey(key string) (*storage.Entity, error)
}

// MetricCollector - сборщик статистики.
type MetricCollector struct {
	storage StorageInterface
}

// NewMetricCollector - конструктор.
func NewMetricCollector(storage StorageInterface) *MetricCollector {
	return &MetricCollector{storage}
}

// Save - собирает статистику.
func (c *MetricCollector) Save(t metric.Type, n string, v any) error {
	k := metric.GetKey(string(t), n)
	memItem, err := c.storage.FindByKey(k)
	if err != nil {
		return err
	}

	if memItem == nil {
		memItem = &storage.Entity{
			Key:  k,
			Name: n,
			Type: t,
		}
	}

	switch t {
	case metric.TypeCounter:
		var val int64
		switch t := v.(type) {
		case string:
			val, err = strconv.ParseInt(t, 10, 64)
			if err != nil {
				return err
			}
		case int64:
			val = t
		default:
			return errors.New("invalid type")
		}

		if memItem.Value == nil {
			memItem.Value = val
		} else {
			memItem.Value = memItem.Value.(int64) + val
		}
	case metric.TypeGauge:
		switch t := v.(type) {
		case string:
			memItem.Value, err = strconv.ParseFloat(t, 64)
			if err != nil {
				return err
			}
		case float64:
			memItem.Value = t
		default:
			return errors.New("invalid type")
		}
	}

	return c.storage.Save(memItem)
}

// GetAll - возвращает все записи в виде DTO.
func (c *MetricCollector) GetAll() map[string]*storage.Entity {
	return c.storage.GetAll()
}

// FindByKey - находит запись по ключу.
func (c *MetricCollector) FindByKey(key string) (*storage.Entity, error) {
	entity, err := c.storage.FindByKey(key)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}

	return entity, nil
}
