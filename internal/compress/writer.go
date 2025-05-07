package compress

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
)

var encodings = map[Type]string{
	Gzip:    "gzip",
	Deflate: "deflate",
	Br:      "br",
}

// Writer структура для записи сжатых данных.
type Writer struct {
	cmp io.WriteCloser
}

// NewWriter конструктор.
func NewWriter(t Type, bb *bytes.Buffer) (*Writer, error) {
	cmp, err := compressor(t, bb)
	if err != nil {
		return nil, err
	}

	return &Writer{
		cmp: cmp,
	}, nil
}

// Write запись.
func (w *Writer) Write(b []byte) (int, error) {
	return w.cmp.Write(b)
}

// Close закрытие.
func (w *Writer) Close() error {
	return w.cmp.Close()
}

// HTTPWriter структура для обработки сжатия ответа.
type HTTPWriter struct {
	Writer
	writer          http.ResponseWriter
	cmp             io.WriteCloser
	contentEncoding string
}

// NewHTTPWriter конструктор.
func NewHTTPWriter(t Type, w http.ResponseWriter) (*HTTPWriter, error) {
	cmp, _ := compressor(t, w)

	httpWr := HTTPWriter{
		writer:          w,
		cmp:             cmp,
		contentEncoding: encodings[t],
	}

	return &httpWr, nil
}

// Header возвращает заголовоки.
func (c *HTTPWriter) Header() http.Header {
	return c.writer.Header()
}

// Write записывает данные.
func (c *HTTPWriter) Write(p []byte) (int, error) {
	if c.cmp == nil {
		return c.writer.Write(p)
	}
	return c.cmp.Write(p)
}

// WriteHeader устанавливает заголовок ответа.
func (c *HTTPWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.writer.Header().Set("Content-Encoding", c.contentEncoding)
	}
	c.writer.WriteHeader(statusCode)
}

// Close закрывает HTTPWriter и досылает все данные из буфера.
func (c *HTTPWriter) Close() error {
	if c.cmp == nil {
		return nil
	}
	return c.cmp.Close()
}

func compressor(t Type, w io.Writer) (io.WriteCloser, error) {
	switch t {
	case Gzip:
		return gzip.NewWriter(w), nil
	case Deflate:
		return zlib.NewWriter(w), nil
	case Br:
		return brotli.NewWriter(w), nil
	}
	return nil, fmt.Errorf("unsupported compress type: %v", t)
}
