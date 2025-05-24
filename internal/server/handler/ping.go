package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/ktigay/metrics-collector/internal/log"
	"github.com/ktigay/metrics-collector/internal/retry"
	"github.com/ktigay/metrics-collector/internal/server/db"
)

const (
	pingTimeout = 2 * time.Second
)

// PingHandler пинг хендлер.
func PingHandler(w http.ResponseWriter, r *http.Request) {
	if db.MasterDB == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), pingTimeout)
	defer cancel()

	handler := func(_ retry.RetPolicy) error {
		return db.MasterDB.PingContext(ctx)
	}
	if err = retry.Ret(handler); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.AppLogger.Errorf("Failed to connect to MasterDB %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
