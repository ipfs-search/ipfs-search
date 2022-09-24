package extractor

import (
	"errors"
)

var (
	// ErrFileTooLarge is returned when the size of a file is larger than the configured `MaxFileSize`.
	ErrFileTooLarge = errors.New("file too large")
)
