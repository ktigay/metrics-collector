package item

import (
	"errors"
	"github.com/ktigay/metrics-collector/internal/metric"
	"strconv"
	"strings"
)

var (
	ErrorInvalidLength = errors.New("invalid length")
	ErrorInvalidVal    = errors.New("invalid value")
	ErrorInvalidType   = errors.New("invalid type")
	ErrorInvalidName   = errors.New("invalid name")
)

// ParseFromPath - парсит метрику из строки.
func ParseFromPath(path string) (m MetricDTO, err error) {

	p := strings.Split(strings.TrimPrefix(path, "/update/"), "/")
	l := len(p)

	switch {
	case l > 1 && p[1] == "":
		return m, ErrorInvalidName
	case l != 3:
		return m, ErrorInvalidLength
	}

	typeStr, metricName, val := p[0], p[1], p[2]

	mType, err := metric.ResolveType(typeStr)
	if err != nil {
		return m, ErrorInvalidType
	}

	switch mType {
	case metric.TypeGauge:
		if i, err := strconv.ParseFloat(val, 64); err == nil {
			m.FloatValue = i
		} else {
			return m, ErrorInvalidVal
		}
	case metric.TypeCounter:
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			m.IntValue = i
		} else {
			return m, ErrorInvalidVal
		}
	}

	m.Type = mType
	m.Name = metricName

	return m, nil
}
