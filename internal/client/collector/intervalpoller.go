package collector

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// StatPoller сборка стат данных.
type StatPoller interface {
	PollStat(ctx context.Context)
}

// IntervalPoller собирает статистику.
type IntervalPoller struct {
	cl       StatPoller
	interval time.Duration
	logger   *zap.SugaredLogger
}

// PollStat сбор статистики.
func (m *IntervalPoller) PollStat(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	m.cl.PollStat(ctx)
	for {
		select {
		case <-ticker.C:
			m.cl.PollStat(ctx)
		case <-ctx.Done():
			m.logger.Debug("pollStat stopped")
			return
		}
	}
}

// NewIntervalPoller собирает статистику.
func NewIntervalPoller(cl StatPoller, pollInterval time.Duration, logger *zap.SugaredLogger) *IntervalPoller {
	return &IntervalPoller{
		cl:       cl,
		interval: pollInterval,
		logger:   logger,
	}
}
