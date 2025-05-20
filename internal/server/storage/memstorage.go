package storage

import (
	"maps"
	"slices"
	"sync"

	"github.com/ktigay/metrics-collector/internal/log"
	"github.com/ktigay/metrics-collector/internal/metric"
)

// Snapshot интерфейс для чтения/сохранения снимка данных.
type Snapshot interface {
	Read() ([]Entity, error)
	Write([]Entity) error
}

// MemStorage - in-memory хранилище.
type MemStorage struct {
	sm       sync.RWMutex
	Metrics  map[string]Entity
	snapshot Snapshot
}

// NewMemStorage - конструктор.
func NewMemStorage(snapshot Snapshot) (*MemStorage, error) {
	storage := MemStorage{
		snapshot: snapshot,
		Metrics:  make(map[string]Entity),
	}
	if err := storage.restore(); err != nil {
		return nil, err
	}

	return &storage, nil
}

// Save - сохраняет метрику.
func (s *MemStorage) Save(m Entity) error {
	s.sm.Lock()
	defer s.sm.Unlock()

	s.Metrics[m.Key] = m
	return nil
}

// Find - поиск по ключу
func (s *MemStorage) Find(t, n string) (*Entity, error) {
	s.sm.RLock()
	defer s.sm.RUnlock()

	key := metric.Key(t, n)
	entity, ok := s.Metrics[key]
	if !ok {
		return nil, nil
	}

	return &entity, nil
}

// All - вернуть все метрики
func (s *MemStorage) All() []Entity {
	s.sm.RLock()
	defer s.sm.RUnlock()

	all := make([]Entity, 0, len(s.Metrics))
	for _, v := range s.Metrics {
		all = append(all, v)
	}
	return all
}

// Remove удаляет по ключу.
func (s *MemStorage) Remove(t, n string) error {
	s.sm.Lock()
	defer s.sm.Unlock()

	key := metric.Key(t, n)
	delete(s.Metrics, key)
	return nil
}

// Backup бэкап данных.
func (s *MemStorage) Backup() error {
	if s.snapshot == nil {
		return nil
	}

	s.sm.RLock()
	defer s.sm.RUnlock()

	return s.snapshot.Write(
		slices.Collect(maps.Values(s.Metrics)),
	)
}

func (s *MemStorage) restore() error {
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
