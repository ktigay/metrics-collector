package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const (
	serverURL = "http://localhost:8080"
)

func ExampleMetricHandler_CollectHandler() {
	// пост запрос для сохранения значения метрики.
	req, _ := http.NewRequest(http.MethodPost, serverURL+"/update/counter/TestSet91/100", nil)

	contentType := []string{"text/html"}

	req.Header = http.Header{
		"Content-Type": contentType,
		"Accept":       contentType,
	}

	resp, _ := http.DefaultClient.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(b))
}

func ExampleMetricHandler_GetValueHandler() {
	// гет запрос для получения значения метрики.
	req, _ := http.NewRequest(http.MethodGet, serverURL+"/value/counter/TestSet91", nil)

	contentType := []string{"text/html"}

	req.Header = http.Header{
		"Content-Type": contentType,
		"Accept":       contentType,
	}

	resp, _ := http.DefaultClient.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(b))
}

func ExampleMetricHandler_GetAllHandler() {
	// гет запрос для получения списка метрик.
	req, _ := http.NewRequest(http.MethodGet, serverURL+"/", nil)

	contentType := []string{"text/html"}

	req.Header = http.Header{
		"Content-Type": contentType,
		"Accept":       contentType,
	}

	resp, _ := http.DefaultClient.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(b))
}

func ExampleMetricHandler_UpdateJSONHandler() {
	// пост запрос на обновление метрики из json.
	req, _ := http.NewRequest(http.MethodPost, serverURL+"/update/", bytes.NewReader([]byte(`{"id":"TestSet91","type":"counter","delta":15,"value":0}`)))

	contentType := []string{"application/json"}

	req.Header = http.Header{
		"Content-Type": contentType,
		"Accept":       contentType,
	}

	resp, _ := http.DefaultClient.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(b))
}

func ExampleMetricHandler_UpdatesJSONHandler() {
	// пост запрос на обновление метрик из json.
	req, _ := http.NewRequest(http.MethodPost, serverURL+"/updates/", bytes.NewReader([]byte(`
		[{"id":"TestSet91","type":"counter","delta":15,"value":0},{"id":"TestSet92","type":"gauge","delta":0,"value":1000}]
	`)))

	contentType := []string{"application/json"}

	req.Header = http.Header{
		"Content-Type": contentType,
		"Accept":       contentType,
	}

	resp, _ := http.DefaultClient.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(b))
}

func ExampleMetricHandler_GetJSONValueHandler() {
	// пост запрос на получение метрики в виде json.
	req, _ := http.NewRequest(http.MethodPost, serverURL+"/value/", bytes.NewReader([]byte(`{"id":"TestSet91","type":"counter"}`)))

	contentType := []string{"application/json"}

	req.Header = http.Header{
		"Content-Type": contentType,
		"Accept":       contentType,
	}

	resp, _ := http.DefaultClient.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(b))
}
