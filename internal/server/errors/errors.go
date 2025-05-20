package errors

import "errors"

var (
	// ErrWrongType направильный тип метрики.
	ErrWrongType = errors.New("wrong type")
	// ErrWrongValue неправильное значение метрики.
	ErrWrongValue = errors.New("wrong value")
	// ErrInvalidValueType неправильный тип значения.
	ErrInvalidValueType = errors.New("invalid value type")
	// ErrValueNotFound значение не найдено.
	ErrValueNotFound = errors.New("value not found")
)
