package errors

// StatusCodeError ошибочный http код.
type StatusCodeError struct {
	Message    string
	StatusCode int
}

// Error ошибка.
func (e *StatusCodeError) Error() string {
	return e.Message
}

// UnsupportedTypeError тип сжатия не поддерживается.
type UnsupportedTypeError struct {
	Type    string
	Message string
}

// Error ошибка.
func (e *UnsupportedTypeError) Error() string {
	return e.Message
}
