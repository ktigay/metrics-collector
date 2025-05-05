package storage

// MemStorage - in-memory хранилище.
type MemStorage struct {
	Metrics map[string]Entity
}

// NewMemStorage - конструктор.
func NewMemStorage(ms []Entity) *MemStorage {
	var mm = make(map[string]Entity)
	for _, m := range ms {
		mm[m.Key] = m
	}
	return &MemStorage{mm}
}

// Save - сохраняет метрику.
func (s *MemStorage) Save(m Entity) error {
	s.Metrics[m.Key] = m
	return nil
}

// FindByKey - поиск по ключу
func (s *MemStorage) FindByKey(key string) (*Entity, error) {
	entity, ok := s.Metrics[key]
	if !ok {
		return nil, nil
	}

	return &entity, nil
}

// GetAll - вернуть все метрики
func (s *MemStorage) GetAll() []Entity {
	var all = make([]Entity, 0)
	for _, v := range s.Metrics {
		all = append(all, v)
	}
	return all
}

func (s *MemStorage) RemoveByKey(key string) error {
	delete(s.Metrics, key)
	return nil
}
