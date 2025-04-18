package storage

import "github.com/ktigay/metrics-collector/internal/metric"

// MemStorageEntity - структура для сохранения в in-memory storage.
type MemStorageEntity struct {
	Key        string
	Type       metric.Type
	Name       string
	IntValue   int64
	FloatValue float64
}
