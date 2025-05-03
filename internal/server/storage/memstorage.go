package storage

// MemStorage - in-memory хранилище.
type MemStorage struct {
	Metrics map[string]*Entity
}

// NewMemStorage - конструктор.
func NewMemStorage(m *map[string]*Entity) *MemStorage {
	if m == nil {
		mp := make(map[string]*Entity)
		m = &mp
	}
	return &MemStorage{*m}
}

// Save - сохраняет метрику.
func (s *MemStorage) Save(m *Entity) error {
	s.Metrics[m.Key] = m
	return nil
}

// FindByKey - поиск по ключу
func (s *MemStorage) FindByKey(key string) (*Entity, error) {
	entity, ok := s.Metrics[key]
	if !ok {
		return nil, nil
	}

	return entity, nil
}

// GetAll - вернуть все метрики
func (s *MemStorage) GetAll() *map[string]*Entity {
	return &s.Metrics
}
