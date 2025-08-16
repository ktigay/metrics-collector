package compress

import (
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/andybalholm/brotli"

	cerr "github.com/ktigay/metrics-collector/internal/compress/errors"
)

// HTTPWriter структура для обработки сжатия ответа.
type HTTPWriter struct {
	writer     http.ResponseWriter
	compressor io.WriteCloser
}

// NewHTTPWriter конструктор.
func NewHTTPWriter(t Type, w http.ResponseWriter) (*HTTPWriter, error) {
	cmp, _ := compressor(t, w)
	if cmp != nil {
		w.Header().Set("Content-Encoding", fmt.Sprint(t))
	}

	httpWr := HTTPWriter{
		writer:     w,
		compressor: cmp,
	}

	return &httpWr, nil
}

// Header возвращает заголовоки.
func (c *HTTPWriter) Header() http.Header {
	return c.writer.Header()
}

// Write записывает данные.
func (c *HTTPWriter) Write(p []byte) (int, error) {
	if c.compressor == nil {
		return c.writer.Write(p)
	}
	return c.compressor.Write(p)
}

// WriteHeader устанавливает заголовок ответа.
func (c *HTTPWriter) WriteHeader(statusCode int) {
	c.writer.WriteHeader(statusCode)
}

// Close закрывает HTTPWriter и досылает все данные из буфера.
func (c *HTTPWriter) Close() error {
	if c.compressor == nil {
		return nil
	}
	return c.compressor.Close()
}

// WriteCloser обертка над потоком сжатия данных.
type WriteCloser struct {
	body []byte
	comp io.WriteCloser
}

// Write записывает данные в поток.
func (b *WriteCloser) Write(p []byte) (n int, err error) {
	b.body = append(b.body, p...)
	return b.comp.Write(p)
}

// Close закрывает поток.
func (b *WriteCloser) Close() error {
	return b.comp.Close()
}

// RawBody оригинальное тело запроса.
func (b *WriteCloser) RawBody() []byte {
	return b.body
}

// NewWriteCloser конструктор.
func NewWriteCloser(t Type, w io.Writer) (*WriteCloser, error) {
	comp, err := compressor(t, w)
	if err != nil {
		return nil, err
	}

	return &WriteCloser{
		comp: comp,
	}, nil
}

// JSON json структуры.
func JSON(w io.WriteCloser, i any, logger Logger) error {
	var err error

	defer func() {
		if e := w.Close(); e != nil {
			logger.Errorf("JSON compressor close error: %v", e)
		}
	}()
	if err = json.NewEncoder(w).Encode(i); err != nil {
		return err
	}
	return nil
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
	return nil, &cerr.UnsupportedTypeError{
		Type:    t.String(),
		Message: fmt.Sprintf("unsupported compress type: %v", t),
	}
}
