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

// GetKey - возвращает ключ по типу и наименованию метрики.
func GetKey(mType string, mName string) string {
	return string(mType) + ":" + mName
}
