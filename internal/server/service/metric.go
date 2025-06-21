// Package service сервисы.
package service

import (
	"context"
	"time"

	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/retry"
	e "github.com/ktigay/metrics-collector/internal/server/errors"
	"github.com/ktigay/metrics-collector/internal/server/repository"
	"go.uber.org/zap"
)

// MetricRepository интерфейс хранилища.
type MetricRepository interface {
	Upsert(ctx context.Context, m repository.MetricEntity) error
	Find(ctx context.Context, t, n string) (*repository.MetricEntity, error)
	Remove(ctx context.Context, t, n string) error
	All(ctx context.Context) ([]repository.MetricEntity, error)
}

// BackupRepository интерфейс для работы со снапшотами.
type BackupRepository interface {
	Backup(ctx context.Context) error
	Restore(ctx context.Context) error
}

// BatchMetricRepository интерфейс для работы с батчами.
type BatchMetricRepository interface {
	UpsertAll(ctx context.Context, mt []repository.MetricEntity) error
}

// MetricCollector сборщик статистики.
type MetricCollector struct {
	repo   MetricRepository
	logger *zap.SugaredLogger
}

// NewMetricCollector конструктор.
func NewMetricCollector(repo MetricRepository, logger *zap.SugaredLogger) *MetricCollector {
	return &MetricCollector{
		repo:   repo,
		logger: logger,
	}
}

// Save собирает статистику.
func (c *MetricCollector) Save(ctx context.Context, mt metric.Metrics) error {
	var (
		t       metric.Type
		memItem repository.MetricEntity
		err     error
	)

	if t, err = metric.ResolveType(mt.Type); err != nil {
		return e.ErrWrongType
	}

	memItem = repository.MetricEntity{
		Key:   mt.Key(),
		Name:  mt.ID,
		Type:  t,
		Delta: mt.GetDelta(),
		Value: mt.GetValue(),
	}

	return c.repo.Upsert(ctx, memItem)
}

// All возвращает все записи.
func (c *MetricCollector) All(ctx context.Context) ([]repository.MetricEntity, error) {
	return c.repo.All(ctx)
}

// Find находит запись по ключу.
func (c *MetricCollector) Find(ctx context.Context, t, n string) (*metric.Metrics, error) {
	var (
		entity *repository.MetricEntity
		err    error
	)
	if _, err = metric.ResolveType(t); err != nil {
		return nil, e.ErrWrongType
	}

	if entity, err = c.repo.Find(ctx, t, n); err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, e.ErrValueNotFound
	}

	m := entity.ToMetrics()
	return &m, nil
}

// Remove удаление записи по ключу.
func (c *MetricCollector) Remove(ctx context.Context, t, n string) error {
	if _, err := metric.ResolveType(t); err != nil {
		return e.ErrWrongType
	}
	return c.repo.Remove(ctx, t, n)
}

// Backup бэкап данных.
func (c *MetricCollector) Backup(mainCtx, exitCtx context.Context, storeInterval int) error {
	var repo BackupRepository
	switch t := c.repo.(type) {
	case BackupRepository:
		repo = t
	}

	if repo == nil {
		c.logger.Debug("repo not supported backups")
		return nil
	}

	ticker := time.NewTicker(time.Duration(storeInterval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.logger.Debug("saveSnapshot saving metrics started")
			if err := repo.Backup(mainCtx); err != nil {
				return err
			}
			c.logger.Debug("saveSnapshot saving metrics finished")
		case <-exitCtx.Done():
			ticker.Stop()
			if err := repo.Backup(mainCtx); err != nil {
				return err
			}
			c.logger.Debug("saveSnapshot shutting down")
			return nil
		}
	}
}

// Restore восстановление данных.
func (c *MetricCollector) Restore(ctx context.Context) error {
	switch t := c.repo.(type) {
	case BackupRepository:
		var err error
		retry.Ret(func(policy retry.Policy) bool {
			err = t.Restore(ctx)
			c.logger.Debugf("try to restore repo retries %d, prev %v", policy.RetIndex()+1, err)
			return err == nil
		})
		return err
	}

	c.logger.Debug("repo not supported restores")
	return nil
}

// SaveAll сохраняет батч.
func (c *MetricCollector) SaveAll(ctx context.Context, mt []metric.Metrics) error {
	var err error
	entities := make([]repository.MetricEntity, 0, len(mt))

	for _, m := range mt {
		var t metric.Type
		if t, err = metric.ResolveType(m.Type); err != nil {
			return e.ErrWrongType
		}

		en := repository.MetricEntity{
			Key:   metric.Key(m.Type, m.ID),
			Name:  m.ID,
			Type:  t,
			Delta: m.GetDelta(),
			Value: m.GetValue(),
		}
		entities = append(entities, en)
	}

	switch t := c.repo.(type) {
	case BatchMetricRepository:
		return t.UpsertAll(ctx, entities)
	default:
		for _, en := range entities {
			if err = c.repo.Upsert(ctx, en); err != nil {
				return err
			}
		}
		return nil
	}
}
