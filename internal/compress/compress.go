package compress

import "strings"

// Type тип сжатия.
type Type string

const (
	Gzip    Type = "gzip"
	Deflate Type = "deflate"
	Br      Type = "br"
)

// TypeFromString поиск типа сжатия из строки.
func TypeFromString(str string) Type {
	if str == "" {
		return ""
	}

	a := strings.Split(str, ",")
	for _, v := range a {
		v = strings.TrimSpace(v)
		switch v {
		case string(Gzip):
			return Gzip
		case string(Br):
			return Br
		case string(Deflate):
			return Deflate
		}
	}
	return ""
}
