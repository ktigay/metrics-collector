package http

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/ktigay/metrics-collector/internal/compress"
)

// CompressWriter структура для обработки сжатия ответа.
type CompressWriter struct {
	writer          http.ResponseWriter
	cmp             io.WriteCloser
	contentEncoding string
}

// Header возвращает заголовоки.
func (c *CompressWriter) Header() http.Header {
	return c.writer.Header()
}

// Write записывает данные.
func (c *CompressWriter) Write(p []byte) (int, error) {
	if c.cmp == nil {
		return c.writer.Write(p)
	}
	return c.cmp.Write(p)
}

// WriteHeader устанавливает заголовок ответа.
func (c *CompressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.writer.Header().Set("Content-Encoding", c.contentEncoding)
	}
	c.writer.WriteHeader(statusCode)
}

// Close закрывает Writer и досылает все данные из буфера.
func (c *CompressWriter) Close() error {
	if c.cmp == nil {
		return nil
	}
	return c.cmp.Close()
}

// CompressWriterFactory фабрика CompressWriter.
func CompressWriterFactory(t compress.Type, w http.ResponseWriter) *CompressWriter {
	switch t {
	case compress.Gzip:
		return newGzipCompressWriter(w)
	case compress.Deflate:
		return newDeflateCompressWriter(w)
	case compress.Br:
		return newBrotliCompressWriter(w)
	default:
		return &CompressWriter{
			writer: w,
		}
	}
}

func newGzipCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		writer:          w,
		cmp:             gzip.NewWriter(w),
		contentEncoding: "gzip",
	}
}

func newDeflateCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		writer:          w,
		cmp:             zlib.NewWriter(w),
		contentEncoding: "deflate",
	}
}

func newBrotliCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		writer:          w,
		cmp:             brotli.NewWriter(w),
		contentEncoding: "br",
	}
}
