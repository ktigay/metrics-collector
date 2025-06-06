package transport

import (
	"io"
	"net/http"

	"github.com/ktigay/metrics-collector/internal/compress"
	"github.com/ktigay/metrics-collector/internal/metric"
	"go.uber.org/zap"
)

// HTTPClient http транспорт отправки метрик.
type HTTPClient struct {
	url          string
	compressType compress.Type
	logger       *zap.SugaredLogger
}

// NewHTTPClient конструктор.
func NewHTTPClient(url string, logger *zap.SugaredLogger) *HTTPClient {
	return &HTTPClient{
		url:          url,
		compressType: compress.Gzip,
		logger:       logger,
	}
}

// Send отправка одной метрики.
func (h *HTTPClient) Send(body metric.Metrics) ([]byte, error) {
	return h.send(h.url+"/update/", body)
}

// SendBatch отправка батча.
func (h *HTTPClient) SendBatch(body []metric.Metrics) ([]byte, error) {
	return h.send(h.url+"/updates/", body)
}

func (h *HTTPClient) send(url string, body any) ([]byte, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
	)

	if req, err = compress.NewJSONRequest(
		http.MethodPost,
		url,
		h.compressType,
		body,
		h.logger,
	); err != nil {
		return nil, err
	}

	if resp, err = compress.NewClient().Do(req); err != nil {
		return nil, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			h.logger.Error("client.post error", zap.Error(err))
		}
	}()

	return io.ReadAll(resp.Body)
}
