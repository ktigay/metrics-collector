package repository

import (
	"context"
	"maps"
	"slices"
	"sync"

	"github.com/ktigay/metrics-collector/internal/metric"
	"go.uber.org/zap"
)

// MetricSnapshot интерфейс для чтения/сохранения снимка данных.
type MetricSnapshot interface {
	Read() ([]MetricEntity, error)
	Write([]MetricEntity) error
}

// MemMetricRepository in-memory хранилище.
type MemMetricRepository struct {
	sm       sync.RWMutex
	Metrics  map[string]MetricEntity
	snapshot MetricSnapshot
	logger   *zap.SugaredLogger
}

// NewMemRepository конструктор.
func NewMemRepository(snapshot MetricSnapshot, logger *zap.SugaredLogger) (*MemMetricRepository, error) {
	repo := MemMetricRepository{
		snapshot: snapshot,
		Metrics:  make(map[string]MetricEntity),
		logger:   logger,
	}

	return &repo, nil
}

// Upsert сохраняет или обновляет существующую метрику.
func (s *MemMetricRepository) Upsert(_ context.Context, m MetricEntity) error {
	s.sm.Lock()
	defer s.sm.Unlock()

	old, exists := s.Metrics[m.Key]
	if exists {
		old.Value = m.Value
		old.Delta += m.Delta
		s.Metrics[m.Key] = old
	} else {
		s.Metrics[m.Key] = m
	}
	return nil
}

// Find поиск по ключу.
func (s *MemMetricRepository) Find(_ context.Context, t, n string) (*MetricEntity, error) {
	s.sm.RLock()
	defer s.sm.RUnlock()

	key := metric.Key(t, n)
	entity, ok := s.Metrics[key]
	if !ok {
		return nil, nil
	}

	return &entity, nil
}

// All вернуть все метрики.
func (s *MemMetricRepository) All(_ context.Context) ([]MetricEntity, error) {
	s.sm.RLock()
	defer s.sm.RUnlock()

	all := make([]MetricEntity, 0, len(s.Metrics))
	for _, v := range s.Metrics {
		all = append(all, v)
	}
	return all, nil
}

// Remove удаляет по типу и наименованию.
func (s *MemMetricRepository) Remove(_ context.Context, t, n string) error {
	s.sm.Lock()
	defer s.sm.Unlock()

	key := metric.Key(t, n)
	delete(s.Metrics, key)
	return nil
}

// Backup бэкап данных в снапшот.
func (s *MemMetricRepository) Backup(_ context.Context) error {
	if s.snapshot == nil {
		return nil
	}

	s.sm.RLock()
	defer s.sm.RUnlock()

	return s.snapshot.Write(
		slices.Collect(maps.Values(s.Metrics)),
	)
}

// Restore восстановление данных из снапшота.
func (s *MemMetricRepository) Restore(_ context.Context) error {
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

	s.logger.Debugf("repository.restore restored len=%d", len(data))
	return nil
}
