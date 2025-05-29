// Package repository Репозитории.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ktigay/metrics-collector/internal/metric"
)

var (
	upsertQuery = `
	insert into metrics (type, name, delta, value)
	values ($1, $2, $3, $4)
	ON CONFLICT ON CONSTRAINT type_name_uidx DO UPDATE
		SET delta      = metrics.delta + EXCLUDED.delta,
			value      = EXCLUDED.value,
			updated_at = now()
	`
	replaceQuery = `
	insert into metrics (type, name, delta, value)
	values ($1, $2, $3, $4)
	ON CONFLICT ON CONSTRAINT type_name_uidx DO UPDATE
		SET delta      = metrics.delta,
			value      = EXCLUDED.value,
			updated_at = now()
	`
	findQuery = `
	SELECT 
			"type", "name", "delta", "value" 
		FROM metrics 
		WHERE "type" = $1
		AND "name" = $2
	`
	removeQuery    = `delete from metrics where type = $1 and name = $2`
	selectAllQuery = `select type, name, delta, value from metrics`
)

const (
	timeout = 1 * time.Second
)

// DBMetricRepository репозиторий БД.
type DBMetricRepository struct {
	db       *sql.DB
	snapshot MetricSnapshot
}

// NewDBMetricRepository конструктор.
func NewDBMetricRepository(db *sql.DB, snapshot MetricSnapshot) (*DBMetricRepository, error) {
	return &DBMetricRepository{
		db:       db,
		snapshot: snapshot,
	}, nil
}

// Upsert сохраняет или обновляет существующую метрику.
func (dbm *DBMetricRepository) Upsert(ctx context.Context, m MetricEntity) error {
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := dbm.db.ExecContext(c, upsertQuery,
		m.Type, m.Name, m.Delta, m.Value)
	return err
}

// Find поиск по ключу.
func (dbm *DBMetricRepository) Find(ctx context.Context, t, n string) (*MetricEntity, error) {
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	m := MetricEntity{}
	r := dbm.db.QueryRowContext(c, findQuery, t, n)
	if err := r.Scan(&m.Type, &m.Name, &m.Delta, &m.Value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	m.Key = metric.Key(fmt.Sprint(m.Type), n)

	return &m, nil
}

// Remove удаляет по типу и наименованию.
func (dbm *DBMetricRepository) Remove(ctx context.Context, t, n string) error {
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if _, err := dbm.db.ExecContext(c, removeQuery, t, n); err != nil {
		return err
	}
	return nil
}

// All вернуть все метрики.
func (dbm *DBMetricRepository) All(ctx context.Context) ([]MetricEntity, error) {
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	rows, err := dbm.db.QueryContext(c, selectAllQuery)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	entities := make([]MetricEntity, 0)
	for rows.Next() {
		m := MetricEntity{}
		if err = rows.Scan(&m.Type, &m.Name, &m.Delta, &m.Value); err != nil {
			return nil, err
		}
		m.Key = metric.Key(fmt.Sprint(m.Type), m.Name)

		entities = append(entities, m)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return entities, nil
}

// Backup бэкап данных в снапшот.
func (dbm *DBMetricRepository) Backup(ctx context.Context) error {
	if dbm.snapshot == nil {
		return nil
	}

	var (
		entities []MetricEntity
		err      error
	)
	if entities, err = dbm.All(ctx); err != nil {
		return err
	}

	return dbm.snapshot.Write(entities)
}

// Restore восстановление данных из снапшота.
func (dbm *DBMetricRepository) Restore(ctx context.Context) error {
	if dbm.snapshot == nil {
		return nil
	}

	var (
		data []MetricEntity
		err  error
	)

	if data, err = dbm.snapshot.Read(); err != nil {
		return err
	}

	return dbm.batch(ctx, replaceQuery, data)
}

// UpsertAll сохраняет батч.
func (dbm *DBMetricRepository) UpsertAll(ctx context.Context, mt []MetricEntity) error {
	return dbm.batch(ctx, upsertQuery, mt)
}

func (dbm *DBMetricRepository) batch(ctx context.Context, query string, mt []MetricEntity) error {
	var (
		err  error
		tx   *sql.Tx
		stmt *sql.Stmt
	)

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	tx, err = dbm.db.BeginTx(c, nil)
	if err != nil {
		return err
	}
	txOK := false
	defer func() {
		if !txOK {
			_ = tx.Rollback()
		}
	}()

	if stmt, err = tx.Prepare(query); err != nil {
		return err
	}
	for _, m := range mt {
		if _, err = stmt.Exec(m.Type, m.Name, m.Delta, m.Value); err != nil {
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	txOK = true
	return nil
}
