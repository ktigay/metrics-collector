package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/ktigay/metrics-collector/internal/retry"
	"go.uber.org/zap"
)

const (
	pingTimeout = 2 * time.Second
)

// PingHandler структура для обработки ping.
type PingHandler struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

// NewPingHandler конструктор.
func NewPingHandler(db *sql.DB, logger *zap.SugaredLogger) *PingHandler {
	return &PingHandler{
		db:     db,
		logger: logger,
	}
}

// Ping пинг хендлер.
func (p *PingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if p.db == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), pingTimeout)
	defer cancel()

	retry.Ret(func(_ retry.Policy) bool {
		err = p.db.PingContext(ctx)
		return err == nil
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.logger.Errorf("Failed to connect to DB %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
