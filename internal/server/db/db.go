package db

import (
	"database/sql"
	"github.com/ktigay/metrics-collector/internal/log"
)

// MasterDB инстанс БД.
var MasterDB *sql.DB

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
