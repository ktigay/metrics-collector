package collector

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/ktigay/metrics-collector/internal/metric"
)

// StatGetter сборка стат данных.
//
//go:generate mockgen -destination=./mocks/mock_statgetter.go -package=mocks github.com/ktigay/metrics-collector/internal/client/collector StatGetter
type StatGetter interface {
	GetStat() ([]metric.Metrics, error)
}

// IntervalPoller собирает статистику.
type IntervalPoller struct {
	source   StatGetter
	logger   *zap.SugaredLogger
	interval time.Duration
}

// PollStat сбор статистики.
func (m *IntervalPoller) PollStat(ctx context.Context, ch chan<- []metric.Metrics) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.logger.Debug("pollStat collect")

			metrics, err := m.source.GetStat()
			if err != nil {
				m.logger.Warnw("failed to get stat", "error", err)
				continue
			}
			ch <- metrics
		case <-ctx.Done():
			m.logger.Debug("pollStat done")
			return
		}
	}
}

// NewIntervalPoller собирает статистику.
func NewIntervalPoller(source StatGetter, pollInterval time.Duration, logger *zap.SugaredLogger) *IntervalPoller {
	return &IntervalPoller{
		source:   source,
		interval: pollInterval,
		logger:   logger,
	}
}
