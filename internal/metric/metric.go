package metric

import (
	"fmt"
	"runtime"
)

// Type - тип метрики.
type Type string

func (t Type) String() string {
	return string(t)
}

// GaugeMetric - тип собираемых метрик.
type GaugeMetric string

const (
	TypeGauge   Type = "gauge"
	TypeCounter Type = "counter"

	Alloc         GaugeMetric = "Alloc"
	BuckHashSys   GaugeMetric = "BuckHashSys"
	Frees         GaugeMetric = "Frees"
	GCCPUFraction GaugeMetric = "GCCPUFraction"
	GCSys         GaugeMetric = "GCSys"
	HeapAlloc     GaugeMetric = "HeapAlloc"
	HeapIdle      GaugeMetric = "HeapIdle"
	HeapInuse     GaugeMetric = "HeapInuse"
	HeapObjects   GaugeMetric = "HeapObjects"
	HeapReleased  GaugeMetric = "HeapReleased"
	HeapSys       GaugeMetric = "HeapSys"
	LastGC        GaugeMetric = "LastGC"
	Lookups       GaugeMetric = "Lookups"
	MCacheInuse   GaugeMetric = "MCacheInuse"
	MCacheSys     GaugeMetric = "MCacheSys"
	MSpanInuse    GaugeMetric = "MSpanInuse"
	MSpanSys      GaugeMetric = "MSpanSys"
	Mallocs       GaugeMetric = "Mallocs"
	NextGC        GaugeMetric = "NextGC"
	NumForcedGC   GaugeMetric = "NumForcedGC"
	NumGC         GaugeMetric = "NumGC"
	OtherSys      GaugeMetric = "OtherSys"
	PauseTotalNs  GaugeMetric = "PauseTotalNs"
	StackInuse    GaugeMetric = "StackInuse"
	StackSys      GaugeMetric = "StackSys"
	Sys           GaugeMetric = "Sys"
	TotalAlloc    GaugeMetric = "TotalAlloc"

	RandomValue string = "RandomValue"
	PollCount   string = "PollCount"
)

// String - название метрики в строку.
func (m GaugeMetric) String() string {
	return string(m)
}

// Metrics структура для обновления метрик.
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// ValueByType значение метрики в зависимости от типа.
func (m *Metrics) ValueByType() any {
	switch m.MType {
	case string(TypeGauge):
		if m.Value != nil {
			return *m.Value
		}
	case string(TypeCounter):
		if m.Delta != nil {
			return *m.Delta
		}
	}
	return nil
}

// Key ключ метрики.
func (m *Metrics) Key() string {
	return Key(m.MType, m.ID)
}

// ResolveType - получает из строки тип.
func ResolveType(s string) (m Type, err error) {
	switch s {
	case string(TypeGauge):
		return TypeGauge, nil
	case string(TypeCounter):
		return TypeCounter, nil
	default:
		return m, fmt.Errorf("unknown metric type: %s", s)
	}
}

// Key - возвращает ключ по типу и наименованию метрики.
func Key(mType string, mName string) string {
	return fmt.Sprintf("%s:%s", mType, mName)
}

// MapGaugeFromMemStats - преобразует метрики из runtime в map.
func MapGaugeFromMemStats(m runtime.MemStats) map[GaugeMetric]float64 {
	return map[GaugeMetric]float64{
		Alloc:         float64(m.Alloc),
		BuckHashSys:   float64(m.BuckHashSys),
		Frees:         float64(m.Frees),
		GCCPUFraction: m.GCCPUFraction,
		GCSys:         float64(m.GCSys),
		HeapAlloc:     float64(m.HeapAlloc),
		HeapIdle:      float64(m.HeapIdle),
		HeapInuse:     float64(m.HeapInuse),
		HeapObjects:   float64(m.HeapObjects),
		HeapReleased:  float64(m.HeapReleased),
		HeapSys:       float64(m.HeapSys),
		LastGC:        float64(m.LastGC),
		Lookups:       float64(m.Lookups),
		MCacheInuse:   float64(m.MCacheInuse),
		MCacheSys:     float64(m.MCacheSys),
		MSpanInuse:    float64(m.MSpanInuse),
		MSpanSys:      float64(m.MSpanSys),
		Mallocs:       float64(m.Mallocs),
		NextGC:        float64(m.NextGC),
		NumForcedGC:   float64(m.NumForcedGC),
		NumGC:         float64(m.NumGC),
		OtherSys:      float64(m.OtherSys),
		PauseTotalNs:  float64(m.PauseTotalNs),
		StackInuse:    float64(m.StackInuse),
		StackSys:      float64(m.StackSys),
		Sys:           float64(m.Sys),
		TotalAlloc:    float64(m.TotalAlloc),
	}
}
