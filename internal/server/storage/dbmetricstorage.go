package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ktigay/metrics-collector/internal/metric"
)

var (
	upsertSQL = `
	insert into metrics (type, name, delta, value)
	values ($1, $2, $3, $4)
	ON CONFLICT ON CONSTRAINT type_name_uidx DO UPDATE
		SET delta      = metrics.delta + EXCLUDED.delta,
			value      = EXCLUDED.value,
			updated_at = now()
	`
	findSQL = `
	SELECT 
			"type", "name", "delta", "value" 
		FROM metrics 
		WHERE "type" = $1
		AND "name" = $2
	`
	removeSQL    = `delete from metrics where type = $1 and name = $2`
	selectAllSQL = `select type, name, delta, value from metrics`
)

// DBMetricStorage репозиторий БД.
type DBMetricStorage struct {
	db       *sql.DB
	snapshot MetricSnapshot
}

func NewDBMetricStorage(db *sql.DB, snapshot MetricSnapshot) (*DBMetricStorage, error) {
	return &DBMetricStorage{
		db:       db,
		snapshot: snapshot,
	}, nil
}

// Upsert - сохраняет или обновляет существующую метрику.
func (dbm *DBMetricStorage) Upsert(m MetricEntity) error {
	_, err := dbm.db.Exec(upsertSQL,
		m.Type, m.Name, m.Delta, m.Value)
	return err
}

// Find - поиск по ключу.
func (dbm *DBMetricStorage) Find(t, n string) (*MetricEntity, error) {
	m := MetricEntity{}
	r := dbm.db.QueryRow(findSQL, t, n)
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
func (dbm *DBMetricStorage) Remove(t, n string) error {
	if _, err := dbm.db.Exec(removeSQL, t, n); err != nil {
		return err
	}
	return nil
}

// All - вернуть все метрики.
func (dbm *DBMetricStorage) All() ([]MetricEntity, error) {
	rows, err := dbm.db.Query(selectAllSQL)
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
func (dbm *DBMetricStorage) Backup() error {
	if dbm.snapshot == nil {
		return nil
	}

	var (
		entities []MetricEntity
		err      error
	)
	if entities, err = dbm.All(); err != nil {
		return err
	}

	return dbm.snapshot.Write(entities)
}

// Restore восстановление данных из снапшота.
func (dbm *DBMetricStorage) Restore() error {
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

	return dbm.UpsertAll(data)
}

// UpsertAll сохраняет батч.
func (dbm *DBMetricStorage) UpsertAll(mt []MetricEntity) error {
	var (
		err  error
		tx   *sql.Tx
		stmt *sql.Stmt
	)

	tx, err = dbm.db.Begin()
	if err != nil {
		return err
	}
	txOK := false
	defer func() {
		if !txOK {
			_ = tx.Rollback()
		}
	}()

	if stmt, err = tx.Prepare(upsertSQL); err != nil {
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
