package types

import (
	"fmt"
)

// TODO: Remove this.

// WrappedError is a generic error type wrapping underlying errors.
type WrappedError struct {
	Err error
	Msg string
}

// Unwrap returns the underlying error.
func (e *WrappedError) Unwrap() error {
	return e.Err
}

// Error returns the message of the current error.
func (e *WrappedError) Error() string {
	return e.Msg
}

// ProviderError is a WrappedError containing information on a specific provider.
type ProviderError struct {
	WrappedError
	Provider Provider
}

// NewProviderErrorf creates a new error based on a provider, a message and a formatting string,
// optionally wrapping an underlying error.
func NewProviderErrorf(wrapping error, p Provider, format string, a ...interface{}) error {
	return &ProviderError{
		WrappedError{
			Err: wrapping,
			Msg: fmt.Sprintf(format, a...),
		},
		p,
	}
}
