package collector

import (
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/server/metric/item"
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
func (c *MetricCollector) Save(m item.MetricDTO) error {
	k := m.GetKey()
	memItem, err := c.storage.FindByKey(k)
	if err != nil {
		return err
	}

	if memItem == nil {
		memItem = &storage.Entity{}
		memItem.Key = m.GetKey()
		memItem.Name = m.Name
		memItem.Type = m.Type
	}

	switch m.Type {
	case metric.TypeCounter:
		if v, err := strconv.ParseInt(m.Value, 10, 64); err == nil {
			if memItem.Value == nil {
				memItem.Value = v
			} else {
				memItem.Value = memItem.Value.(int64) + v
			}
		} else {
			return err
		}
	case metric.TypeGauge:
		if v, err := strconv.ParseFloat(m.Value, 64); err == nil {
			memItem.Value = v
		} else {
			return err
		}
	}

	return c.storage.Save(memItem)
}

// GetAll - возвращает все записи в виде DTO.
func (c *MetricCollector) GetAll() map[string]item.MetricDTO {
	entities := c.storage.GetAll()
	dtoMap := make(map[string]item.MetricDTO, len(entities))
	for key, m := range entities {
		dtoMap[key] = item.MetricDTO{
			Name:  m.Name,
			Type:  m.Type,
			Value: m.GetValueAsString(),
		}
	}

	return dtoMap
}

// FindByKey - находит запись по ключу.
func (c *MetricCollector) FindByKey(key string) (*item.MetricDTO, error) {
	entity, err := c.storage.FindByKey(key)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}

	return &item.MetricDTO{
		Name:  entity.Name,
		Type:  entity.Type,
		Value: entity.GetValueAsString(),
	}, nil
}
