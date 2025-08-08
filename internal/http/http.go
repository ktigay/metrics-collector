// Package http Writer для логирования данных.
package http

import (
	"crypto/sha256"
	"fmt"
	"net/http"
)

const (
	// HashSHA256Header имя хедера HashSHA256.
	HashSHA256Header = "HashSHA256"
)

type (
	// ResponseData статистика по ответу.
	ResponseData struct {
		Status int
		Size   int
		Body   []byte
		Err    error
	}

	// Writer структура для вывода данных.
	Writer struct {
		http.ResponseWriter
		responseData *ResponseData
		hashKey      string
	}
)

// NewWriter конструктор.
func NewWriter(w http.ResponseWriter, d *ResponseData, hashKey string) *Writer {
	return &Writer{
		ResponseWriter: w,
		responseData:   d,
		hashKey:        hashKey,
	}
}

// Write запись ответа.
func (w *Writer) Write(b []byte) (int, error) {
	w.responseData.Body = append(w.responseData.Body, b...)
	return len(b), nil
}

// WriteHeader устанавливает статус.
func (w *Writer) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.responseData.Status = statusCode
}

// ResponseData данные по ответу.
func (w *Writer) ResponseData() *ResponseData {
	return w.responseData
}

// WithWriter оборачивает исходный src [http.ResponseWriter] в [http.ResponseWriter] возращенный callback функцией.
func (w *Writer) WithWriter(callback func(src http.ResponseWriter) http.ResponseWriter) {
	w.ResponseWriter = callback(w.ResponseWriter)
}

// Flush отправляет клиенту все буферизованные данные.
func (w *Writer) Flush() {
	rw := w.ResponseWriter

	if w.hashKey != "" {
		srvCheckSum := sha256.Sum256(append(w.responseData.Body, w.hashKey...))
		rw.Header().Set(HashSHA256Header, fmt.Sprintf("%x", srvCheckSum))
	}

	if w.responseData.Status == 0 {
		w.WriteHeader(http.StatusOK)
	}

	size, err := rw.Write(w.responseData.Body)
	if err != nil {
		w.responseData.Err = err
	}
	w.responseData.Size = size

	if fl, ok := w.ResponseWriter.(http.Flusher); ok {
		fl.Flush()
	}
}
