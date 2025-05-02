package http

import (
	"compress/gzip"
	"net/http"
)

type CompressWriter struct {
	w               http.ResponseWriter
	zw              *gzip.Writer
	contentEncoding string
}

func (c *CompressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *CompressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *CompressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", c.contentEncoding)
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает Writer и досылает все данные из буфера.
func (c *CompressWriter) Close() error {
	return c.zw.Close()
}

func CompressWriterFactory(t string, w http.ResponseWriter) (http.ResponseWriter, func() error) {
	switch t {
	case "gzip":
		gw := newGzipCompressWriter(w)
		return gw, func() error {
			return gw.Close()
		}
	default:
		return w, func() error {
			return nil
		}
	}
}

func newGzipCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		w:               w,
		zw:              gzip.NewWriter(w),
		contentEncoding: "gzip",
	}
}
