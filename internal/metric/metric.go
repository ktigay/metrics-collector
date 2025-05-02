package metric

import "fmt"

// Type - тип метрики.
type Type string

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
	return mType + ":" + mName
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
