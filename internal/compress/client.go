package compress

import (
	"fmt"
	"io"
	"net/http"

	e "github.com/ktigay/metrics-collector/internal/compress/errors"
)

// NewJSONRequest запрос.
func NewJSONRequest(method, url string, t Type, body any, logger Logger) (*http.Request, error) {
	var (
		err error
		req *http.Request
	)

	r, w := io.Pipe()

	go func() {
		jsonErr := JSON(t, w, body, logger)
		if ec := w.CloseWithError(jsonErr); ec != nil {
			logger.Errorf("w.CloseWithError error: %v", ec)
		}
	}()

	if req, err = http.NewRequest(method, url, r); err != nil {
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

	return req, nil
}

// Client клиент для отправки сжатых запросов.
type Client struct {
	http.Client
}

// NewClient конструктор.
func NewClient() *Client {
	return &Client{}
}

// Do выполнить запрос.
func (c Client) Do(req *http.Request) (*http.Response, error) {
	var (
		err  error
		resp *http.Response
		rc   io.ReadCloser
	)

	if resp, err = c.Client.Do(req); err != nil {
		return nil, err
	}
	if resp.StatusCode > 300 || resp.StatusCode < 200 {
		return nil, &e.StatusCodeError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("status code is not OK %d", resp.StatusCode),
		}
	}

	enc := resp.Header.Get("Content-Encoding")
	if ceAlg := TypeFromString(enc); ceAlg != "" {
		if rc, err = ReaderFactory(ceAlg, resp.Body); err != nil {
			return nil, err
		}
		resp.Body = rc
	}

	return resp, nil
}
