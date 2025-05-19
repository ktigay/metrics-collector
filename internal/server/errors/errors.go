package errors

import "errors"

var (
	ErrWrongType     = errors.New("wrong type")
	ErrWrongValue    = errors.New("wrong value")
	ErrInvalidType   = errors.New("invalid type")
	ErrTypeNotFound  = errors.New("type not found")
	ErrValueNotFound = errors.New("value not found")
)
