// Package db БД.
package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	"github.com/ktigay/metrics-collector/internal/retry"
)

const connTimeout = 100 * time.Millisecond

var structure = []string{
	`
	DO ' BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = ''metric_type'') THEN
		CREATE TYPE metric_type AS ENUM (''counter'', ''gauge'');
		END IF;
	END '
	`,
	`
	CREATE TABLE IF NOT EXISTS metrics
	(
		guid       UUID                     DEFAULT gen_random_uuid(),
		type       metric_type  NOT NULL,
		name       VARCHAR(255) NOT NULL,
		delta      BIGINT                   DEFAULT 0,
		value      DOUBLE PRECISION         DEFAULT .0,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		PRIMARY KEY (guid),
		CONSTRAINT type_name_uidx UNIQUE (type, name)
	)`,
}

// InitializeDB инициализация соединения к БД.
func InitializeDB(ctx context.Context, driver, dsn string, logger *zap.SugaredLogger) (*sql.DB, error) {
	logger.Debugf("Initializing master database %s", dsn)
	dbPool, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	var pingErr error
	retry.Ret(func(policy retry.Policy) bool {
		ctxt, cancel := context.WithTimeout(ctx, connTimeout)
		defer cancel()

		if pingErr = dbPool.PingContext(ctxt); pingErr != nil {
			logger.Debugf("Attempting to connect to dbPool %s, retries %d, prev %v", dsn, policy.RetIndex()+1, pingErr)
			// скип, если это не ошибка соединения
			var pgErr *pgconn.ConnectError
			if !errors.As(pingErr, &pgErr) {
				return true
			}
		}

		return pingErr == nil
	})
	if pingErr != nil {
		return nil, pingErr
	}

	return dbPool, nil
}

// CreateStructure создает структуру БД.
func CreateStructure(ctx context.Context, dbPool *sql.DB) (err error) {
	c, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	for _, s := range structure {
		if _, err = dbPool.ExecContext(c, s); err != nil {
			return nil
		}
	}

	return nil
}

// CloseDB закрывает коннект к БД.
func CloseDB(db *sql.DB) error {
	if db == nil {
		return nil
	}

	return db.Close()
}
