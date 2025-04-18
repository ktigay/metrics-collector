package collector

import (
	"github.com/ktigay/metrics-collector/internal/server/metric/item"
)

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
	return c.storage.Save(m)
}

func (c *MetricCollector) GetAll() map[string]item.MetricDTO {
	return c.storage.GetAll()
}

func (c *MetricCollector) FindByKey(key string) (*item.MetricDTO, error) {
	return c.storage.FindByKey(key)
}

// StorageInterface - интерфейс хранилища.
type StorageInterface interface {
	Save(m item.MetricDTO) error
	GetAll() map[string]item.MetricDTO
	FindByKey(key string) (*item.MetricDTO, error)
}
