package server

import "net/http"

type (
	ResponseData struct {
		Status int
		Size   int
	}

	ResponseWriter struct {
		http.ResponseWriter
		responseData *ResponseData
	}
)

// NewResponseWriter - конструктор.
func NewResponseWriter(w http.ResponseWriter, d *ResponseData) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		responseData:   d,
	}
}

// Write - запись ответа.
func (r *ResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.Size += size
	return size, err
}

// WriteHeader устанавливает статус.
func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.Status = statusCode
}
