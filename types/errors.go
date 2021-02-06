package types

import (
	"errors"
)

// WrappedError is a generic error type wrapping underlying errors.
type WrappedError struct {
	Err error
	Msg string
}

// Unwrap returns the underlying error.
func (e WrappedError) Unwrap() error {
	return e.Err
}

// Error returns the message of the current error.
func (e WrappedError) Error() string {
	return e.Msg
}

var (
	// ErrInvalidResource is returned when a resource is unsupported or invalid.
	ErrInvalidResource = errors.New("resource invalid")

	// ErrUnsupportedType is returned when the type of a resource is currently unsupported.
	ErrUnsupportedType = WrappedError{ErrInvalidResource, "unsupported type"}
)
