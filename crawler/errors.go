package crawler

import (
	"errors"
)

var (
	ErrUnexpectedType = errors.New("unexpected type")
)

// IsTemporaryErr returns true whenever an underlying error signifies a known temporary outage condition rather than permanent failure.
func IsTemporaryErr(err error) bool {
	// TODO: Implement & test me.
	// TODO: Decide whether, with the retrying GET, we actually want this!
	return false
}
