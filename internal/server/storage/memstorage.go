package storage

import (
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/server/metric/item"
)

// MemStorage - in-memory хранилище.
type MemStorage struct {
	metrics map[string]*MemStorageEntity
}

// NewMemStorage - конструктор.
func NewMemStorage() *MemStorage {
	m := make(map[string]*MemStorageEntity)
	return &MemStorage{m}
}

// Save - сохраняет метрику.
func (s *MemStorage) Save(m item.MetricDTO) error {

	k := m.GetKey()
	memItem := s.metrics[k]

	if memItem == nil {
		memItem = &MemStorageEntity{}
		memItem.Key = m.GetKey()
		memItem.Name = m.Name
		memItem.Type = m.Type

		s.metrics[k] = memItem
	}

	switch m.Type {
	case metric.TypeCounter:
		memItem.IntValue += m.IntValue
	case metric.TypeGauge:
		memItem.FloatValue = m.FloatValue
	}

	return nil
}

// FindByKey - поиск по ключу
func (s *MemStorage) FindByKey(key string) (*item.MetricDTO, error) {

	entity, ok := s.metrics[key]
	if !ok {
		return nil, nil
	}

	return &item.MetricDTO{
		Name:       entity.Name,
		Type:       entity.Type,
		IntValue:   entity.IntValue,
		FloatValue: entity.FloatValue,
	}, nil
}

// GetAll - вернуть все метрики
func (s *MemStorage) GetAll() map[string]item.MetricDTO {
	dtos := make(map[string]item.MetricDTO, len(s.metrics))
	for key, m := range s.metrics {
		dtos[key] = item.MetricDTO{
			Name:       m.Name,
			Type:       m.Type,
			IntValue:   m.IntValue,
			FloatValue: m.FloatValue,
		}
	}

	return dtos
}
