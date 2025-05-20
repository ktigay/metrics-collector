package http

import "net/http"

type (
	// ResponseData статистика по ответу.
	ResponseData struct {
		Status int
		Size   int
	}

	// Writer структура для вывода данных.
	Writer struct {
		http.ResponseWriter
		responseData *ResponseData
	}
)

// NewWriter - конструктор.
func NewWriter(w http.ResponseWriter, d *ResponseData) *Writer {
	return &Writer{
		ResponseWriter: w,
		responseData:   d,
	}
}

// Write - запись ответа.
func (r *Writer) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.Size += size
	return size, err
}

// WriteHeader устанавливает статус.
func (r *Writer) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.Status = statusCode
}
