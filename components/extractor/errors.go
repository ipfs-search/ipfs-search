package extractor

import (
	"errors"
)

var (
	// ErrFileTooLarge is returned when the size of a file is larger than the configured `MaxFileSize`.
	ErrFileTooLarge = errors.New("file too large")

	// ErrUnexpectedResponse is returned upon unexpected responses from the backend.
	ErrUnexpectedResponse = errors.New("unexpected response from backend")

	// ErrRequest is returned on errors performing upstream requests.
	ErrRequest = errors.New("request error")
)
