// Package compress Работа со сжатыми данными.
package compress

import "strings"

// Type тип сжатия.
type Type string

const (
	// Gzip gzip.
	Gzip Type = "gzip"
	// Deflate deflate.
	Deflate Type = "deflate"
	// Br brotli.
	Br Type = "br"
)

// String тип в виде строки.
func (t Type) String() string {
	return string(t)
}

// TypeFromString поиск типа сжатия из строки.
func TypeFromString(str string) Type {
	if str == "" {
		return ""
	}

	a := strings.Split(str, ",")
	for _, v := range a {
		v = strings.TrimSpace(v)
		switch Type(v) {
		case Gzip:
			return Gzip
		case Br:
			return Br
		case Deflate:
			return Deflate
		}
	}
	return ""
}

// Logger логгер.
type Logger interface {
	Errorf(string, ...interface{})
}
