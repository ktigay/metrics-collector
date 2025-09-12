package transport

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/ktigay/metrics-collector/internal/compress"
	"github.com/ktigay/metrics-collector/internal/crypto"
	"github.com/ktigay/metrics-collector/internal/metric"
)

const (
	updatePath  = "/update/"
	updatesPath = "/updates/"

	contentType = "application/json"
)

// HTTPClient http транспорт отправки метрик.
type HTTPClient struct {
	requestFactory *RequestFactory
	logger         *zap.SugaredLogger
}

// NewHTTPClient конструктор.
func NewHTTPClient(requestFactory *RequestFactory, logger *zap.SugaredLogger) *HTTPClient {
	return &HTTPClient{
		requestFactory: requestFactory,
		logger:         logger,
	}
}

// Send отправка одной метрики.
func (h *HTTPClient) Send(body metric.Metrics) ([]byte, error) {
	req, err := h.requestFactory.NewRequest(updatePath, body)
	if err != nil {
		return nil, err
	}
	return h.do(req)
}

// SendBatch отправка батча.
func (h *HTTPClient) SendBatch(body []metric.Metrics) ([]byte, error) {
	req, err := h.requestFactory.NewRequest(updatesPath, body)
	if err != nil {
		return nil, err
	}
	return h.do(req)
}

func (h *HTTPClient) do(request *http.Request) ([]byte, error) {
	var (
		err  error
		resp *http.Response
	)

	if resp, err = compress.NewClient().Do(request); err != nil {
		return nil, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			h.logger.Error("client.post error", zap.Error(err))
		}
	}()

	return io.ReadAll(resp.Body)
}

// RequestFactory структура для создания запроса.
type RequestFactory struct {
	encryptKey *crypto.PublicKey
	method     string
	url        string
	hashKey    string
}

// NewRequestFactory конструктор.
func NewRequestFactory(method, url, hashKey string, encryptKey *crypto.PublicKey) *RequestFactory {
	return &RequestFactory{
		method:     method,
		url:        url,
		hashKey:    hashKey,
		encryptKey: encryptKey,
	}
}

// NewRequest создаёт новый запрос.
func (r *RequestFactory) NewRequest(path string, requestBody any) (*http.Request, error) {
	var (
		err error
		req *http.Request
		b   []byte
	)

	if b, err = json.Marshal(requestBody); err != nil {
		return nil, err
	}

	opts := []compress.Option{
		compress.WithHashKey(r.hashKey),
		compress.WithContentType(contentType),
	}
	if r.encryptKey != nil {
		opts = append(opts, compress.WithWriters(func(w io.Writer) io.Writer {
			jw := crypto.Writer{
				Writer: w,
				Key:    r.encryptKey.Key,
			}
			return &jw
		}))
	}

	if req, err = compress.NewRequest(
		r.method,
		r.url+path,
		b,
		opts...,
	); err != nil {
		return nil, err
	}

	return req, nil
}
