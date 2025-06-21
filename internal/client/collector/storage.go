package collector

import (
	"math/rand/v2"
	"sync"

	"github.com/ktigay/metrics-collector/internal/metric"
	"go.uber.org/zap"
)

// ClientMetrics метрики.
type ClientMetrics map[string]*metric.Metrics

// Storage репозиторий для хранения статистики.
type Storage struct {
	mu      sync.Mutex
	metrics ClientMetrics
	logger  *zap.SugaredLogger
}

// Save сохраняет статистику.
func (s *Storage) Save(metrics []metric.Metrics) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range metrics {
		s.metrics[v.ID] = &v
	}

	s.incCounter()
	s.setRandomValue()
}

func (s *Storage) incCounter() {
	pc, ok := s.metrics[metric.PollCount]
	if !ok {
		delta := int64(0)
		s.metrics[metric.PollCount] = &metric.Metrics{
			ID:    metric.PollCount,
			Type:  string(metric.TypeCounter),
			Delta: &delta,
		}
		pc = s.metrics[metric.PollCount]
	}
	v := *pc.Delta + 1
	if err := pc.SetValueByType(v); err != nil {
		s.logger.Errorf("failed to set value by type: %v", err)
	}
}

func (s *Storage) setRandomValue() {
	rnd, ok := s.metrics[metric.RandomValue]
	if !ok {
		s.metrics[metric.RandomValue] = &metric.Metrics{
			ID:   metric.RandomValue,
			Type: string(metric.TypeGauge),
		}
		rnd = s.metrics[metric.RandomValue]
	}
	if err := rnd.SetValueByType(rand.Float64()); err != nil {
		s.logger.Errorf("failed to set value by type: %v", err)
	}
}

// GetStat возвращает собранную статистику.
func (s *Storage) GetStat() []metric.Metrics {
	s.mu.Lock()
	defer s.mu.Unlock()

	metrics := make([]metric.Metrics, 0, len(s.metrics))
	for _, v := range s.metrics {
		metrics = append(metrics, *v)
	}
	return metrics
}

// NewStorage конструктор.
func NewStorage(logger *zap.SugaredLogger) *Storage {
	return &Storage{
		metrics: make(ClientMetrics),
		logger:  logger,
	}
}
