// Package db БД.
package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/ktigay/metrics-collector/internal/log"
)

// MasterDB инстанс БД.
var MasterDB *sql.DB

var structure = []string{
	`
	DO ' BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = ''metric_type'') THEN
		create type metric_type as enum (''counter'', ''gauge'');
		END IF;
	END '
	`,
	`
	CREATE TABLE IF NOT EXISTS metrics
	(
		guid       uuid                     DEFAULT gen_random_uuid(),
		type       metric_type  NOT NULL,
		name       varchar(255) NOT NULL,
		delta      bigint                   default 0,
		value      double precision         default .0,
		created_at timestamp with time zone default now(),
		updated_at timestamp with time zone default now(),
		PRIMARY KEY (guid),
		CONSTRAINT type_name_uidx UNIQUE (type, name)
	)`,
}

// InitializeMasterDB инициализация соединения к БД.
func InitializeMasterDB(driver, dsn string) error {
	log.AppLogger.Debugf("Initializing master database %s", dsn)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}

	MasterDB = db
	return nil
}

// CreateStructure создает структуру БД.
func CreateStructure(ctx context.Context) (err error) {
	c, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	for _, s := range structure {
		if _, err = MasterDB.ExecContext(c, s); err != nil {
			return nil
		}
	}

	return nil
}

// CloseMasterDB закрывает коннект к БД.
func CloseMasterDB() error {
	if MasterDB == nil {
		return nil
	}

	return MasterDB.Close()
}
