package compress

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap/buffer"

	h "github.com/ktigay/metrics-collector/internal/http"
)

// Header тип заголовок.
type Header string

const (
	xRealIP     Header = "X-Real-Ip"
	contentType Header = "Content-Type"
	accept      Header = "Accept"
)

// Options опции реквеста.
type Options struct {
	compressType Type
	hashKey      string
	headers      map[string]string
	writers      []func(w io.Writer) io.Writer
}

// NewOptions конструктор.
func NewOptions(opt []Option) *Options {
	opts := &Options{
		headers: make(map[string]string),
	}
	for _, o := range opt {
		o(opts)
	}
	return opts
}

// Option функция для установки параметров опций.
type Option func(*Options)

// WithCompressType реквест с [Type].
func WithCompressType(compressType Type) Option {
	return func(opt *Options) {
		opt.compressType = compressType
	}
}

// WithHashKey реквест с hashKey.
func WithHashKey(hashKey string) Option {
	return func(opt *Options) {
		opt.hashKey = hashKey
	}
}

// WithContentType реквест с ContentType.
func WithContentType(c string) Option {
	return func(opt *Options) {
		WithHeader(string(contentType), c)(opt)
		WithHeader(string(accept), c)(opt)
	}
}

// WithWriters реквест с [io.Writer] для пост-обработки тела запроса.
func WithWriters(wf ...func(w io.Writer) io.Writer) Option {
	return func(opt *Options) {
		opt.writers = wf
	}
}

// WithXRealIP реквест с заголовком X-Real-IP.
func WithXRealIP(ip string) Option {
	return func(opt *Options) {
		WithHeader(string(xRealIP), ip)(opt)
	}
}

// WithHeader реквест с заголовком.
func WithHeader(key, value string) Option {
	return func(opt *Options) {
		opt.headers[key] = value
	}
}

// NewRequest запрос.
func NewRequest(method, url string, requestBody []byte, opt ...Option) (*http.Request, error) {
	var (
		comp *WriteCloser
		err  error
		req  *http.Request
		t    Type
	)

	opts := NewOptions(opt)

	if opts.compressType == "" {
		t = Default
	} else {
		t = opts.compressType
	}

	w := buffer.Buffer{}

	if comp, err = NewWriteCloser(t, &w); err != nil {
		return nil, err
	}

	if _, err = comp.Write(requestBody); err != nil {
		return nil, err
	}
	if err = comp.Close(); err != nil {
		return nil, err
	}

	b := w.Bytes()
	w.Reset()
	wr := io.Writer(&w)
	for _, wf := range opts.writers {
		wr = wf(wr)
	}
	if _, err = wr.Write(b); err != nil {
		return nil, err
	}

	if req, err = http.NewRequest(method, url, bytes.NewReader(w.Bytes())); err != nil {
		return nil, err
	}

	enc := []string{fmt.Sprint(t)}

	req.Header = http.Header{
		"Content-Encoding": enc,
		"Accept-Encoding":  enc,
	}
	if len(opts.headers) > 0 {
		for k, v := range opts.headers {
			req.Header.Set(k, v)
		}
	}

	rb := comp.RawBody()
	if opts.hashKey != "" && len(rb) > 0 {
		hash := sha256.Sum256(append(rb, opts.hashKey...))
		req.Header[h.HashSHA256Header] = []string{fmt.Sprintf("%x", hash)}
	}

	return req, nil
}
