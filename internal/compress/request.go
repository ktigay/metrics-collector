package compress

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"net/http"

	"go.uber.org/zap/buffer"

	h "github.com/ktigay/metrics-collector/internal/http"
)

// Options опции реквеста.
type Options struct {
	hashKey string
	logger  Logger
}

// NewOptions конструктор.
func NewOptions(opt []Option) *Options {
	opts := &Options{}
	for _, o := range opt {
		o(opts)
	}
	return opts
}

// Option функция для установки параметров опций.
type Option func(*Options)

// WithHashKey реквест с hashKey.
func WithHashKey(hashKey string) Option {
	return func(opt *Options) {
		opt.hashKey = hashKey
	}
}

// WithLogger реквест с логгером.
func WithLogger(logger Logger) Option {
	return func(opt *Options) {
		opt.logger = logger
	}
}

// NewJSONRequest запрос.
func NewJSONRequest(method, url string, t Type, body any, opt ...Option) (*http.Request, error) {
	var (
		comp *WriteCloser
		err  error
		req  *http.Request
	)

	opts := NewOptions(opt)

	w := buffer.Buffer{}
	if comp, err = NewWriteCloser(t, &w); err != nil {
		return nil, err
	}

	if err = JSON(comp, body, opts.logger); err != nil {
		return nil, err
	}

	if req, err = http.NewRequest(method, url, bytes.NewReader(w.Bytes())); err != nil {
		return nil, err
	}

	contentType := []string{"application/json"}
	enc := []string{fmt.Sprint(t)}

	req.Header = http.Header{
		"Content-Type":     contentType,
		"Accept":           contentType,
		"Content-Encoding": enc,
		"Accept-Encoding":  enc,
	}

	rb := comp.RawBody()
	if opts.hashKey != "" && len(rb) > 0 {
		hash := sha256.Sum256(append(rb, opts.hashKey...))
		req.Header[h.HashSHA256Header] = []string{fmt.Sprintf("%x", hash)}
	}

	return req, nil
}
