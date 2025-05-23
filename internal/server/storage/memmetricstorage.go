package storage

import (
	"maps"
	"slices"
	"sync"

	"github.com/ktigay/metrics-collector/internal/log"
	"github.com/ktigay/metrics-collector/internal/metric"
)

// MetricSnapshot интерфейс для чтения/сохранения снимка данных.
type MetricSnapshot interface {
	Read() ([]MetricEntity, error)
	Write([]MetricEntity) error
}

// MemMetricStorage - in-memory хранилище.
type MemMetricStorage struct {
	sm       sync.RWMutex
	Metrics  map[string]MetricEntity
	snapshot MetricSnapshot
}

// NewMemStorage - конструктор.
func NewMemStorage(snapshot MetricSnapshot) (*MemMetricStorage, error) {
	storage := MemMetricStorage{
		snapshot: snapshot,
		Metrics:  make(map[string]MetricEntity),
	}

	return &storage, nil
}

// Save - сохраняет метрику.
func (s *MemMetricStorage) Save(m MetricEntity) error {
	s.sm.Lock()
	defer s.sm.Unlock()

	s.Metrics[m.Key] = m
	return nil
}

// Find - поиск по ключу.
func (s *MemMetricStorage) Find(t, n string) (*MetricEntity, error) {
	s.sm.RLock()
	defer s.sm.RUnlock()

	key := metric.Key(t, n)
	entity, ok := s.Metrics[key]
	if !ok {
		return nil, nil
	}

	return &entity, nil
}

// All - вернуть все метрики.
func (s *MemMetricStorage) All() ([]MetricEntity, error) {
	s.sm.RLock()
	defer s.sm.RUnlock()

	all := make([]MetricEntity, 0, len(s.Metrics))
	for _, v := range s.Metrics {
		all = append(all, v)
	}
	return all, nil
}

// Remove удаляет по типу и наименованию.
func (s *MemMetricStorage) Remove(t, n string) error {
	s.sm.Lock()
	defer s.sm.Unlock()

	key := metric.Key(t, n)
	delete(s.Metrics, key)
	return nil
}

// Backup бэкап данных.
func (s *MemMetricStorage) Backup() error {
	if s.snapshot == nil {
		return nil
	}

	s.sm.RLock()
	defer s.sm.RUnlock()

	return s.snapshot.Write(
		slices.Collect(maps.Values(s.Metrics)),
	)
}

// Restore восстановление данных.
func (s *MemMetricStorage) Restore() error {
	if s.snapshot == nil {
		return nil
	}

	data, err := s.snapshot.Read()
	if err != nil {
		return err
	}
	for _, m := range data {
		s.Metrics[m.Key] = m
	}

	log.AppLogger.Debugf("storage.restore restored len=%d", len(data))

	return nil
}
