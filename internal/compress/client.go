package compress

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// Request структура для формирования сжатого запроса.
type Request struct {
	compress    Type
	contentType string
}

// NewRequest конструктор.
func NewRequest(compress Type, contentType string) *Request {
	return &Request{
		compress:    compress,
		contentType: contentType,
	}
}

// R собирает реквест.
func (r Request) R(method, url string, body []byte) (*http.Request, error) {
	var (
		err  error
		buff bytes.Buffer
		cw   io.WriteCloser
		req  *http.Request
	)

	if cw, err = NewWriteCloser(r.compress, &buff); err != nil {
		return nil, err
	}
	if _, err = cw.Write(body); err != nil {
		return nil, err
	}
	if err = cw.Close(); err != nil {
		return nil, err
	}

	if req, err = http.NewRequest(method, url, &buff); err != nil {
		return nil, err
	}

	enc := string(r.compress)

	req.Header.Set("Content-Type", r.contentType)
	req.Header.Set("Accept", r.contentType)
	req.Header.Set("Content-Encoding", enc)
	req.Header.Set("Accept-Encoding", enc)

	return req, nil
}

// Client клиент для отправки сжатых запросов.
type Client struct{}

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

	client := http.Client{}
	if resp, err = client.Do(req); err != nil {
		return nil, err
	}
	if resp.StatusCode > 300 || resp.StatusCode < 200 {
		return nil, fmt.Errorf("status code is not OK %d", resp.StatusCode)
	}

	enc := resp.Header.Get("Content-Encoding")
	if ceAlg := TypeFromString(enc); ceAlg != "" {
		rc, err = ReaderFactory(ceAlg, resp.Body)
		if err != nil {
			return nil, err
		}
		resp.Body = rc
	}

	return resp, nil
}
