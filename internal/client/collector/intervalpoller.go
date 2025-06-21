package collector

import (
	"context"
	"time"

	"github.com/ktigay/metrics-collector/internal/metric"
	"go.uber.org/zap"
)

// StatGetter сборка стат данных.
//
//go:generate mockgen -destination=./mocks/mock_statgetter.go -package=mocks github.com/ktigay/metrics-collector/internal/client/collector StatGetter
type StatGetter interface {
	GetStat() []metric.Metrics
}

// IntervalPoller собирает статистику.
type IntervalPoller struct {
	source   StatGetter
	interval time.Duration
	logger   *zap.SugaredLogger
}

// PollStat сбор статистики.
func (m *IntervalPoller) PollStat(ctx context.Context, ch chan<- []metric.Metrics) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ch <- m.source.GetStat()
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
