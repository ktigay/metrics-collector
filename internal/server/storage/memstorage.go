package storage

import (
	"maps"
	"slices"
	"sync"
)

type Snapshot interface {
	Read() ([]Entity, error)
	Write([]Entity) error
}

// MemStorage - in-memory хранилище.
type MemStorage struct {
	sm sync.RWMutex
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

// FindByKey - поиск по ключу
func (s *MemStorage) FindByKey(key string) (*Entity, error) {
	s.sm.RLock()
	defer s.sm.RUnlock()

	entity, ok := s.Metrics[key]
	if !ok {
		return nil, nil
	}

	return &entity, nil
}

// GetAll - вернуть все метрики
func (s *MemStorage) GetAll() []Entity {
	s.sm.RLock()
	defer s.sm.RUnlock()

	var all = make([]Entity, 0)
	for _, v := range s.Metrics {
		all = append(all, v)
	}
	return all
}

// RemoveByKey удаляет по ключу.
func (s *MemStorage) RemoveByKey(key string) error {
	s.sm.Lock()
	defer s.sm.Unlock()

	delete(s.Metrics, key)
	return nil
}

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

	return nil
}
