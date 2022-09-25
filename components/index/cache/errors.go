package cache

// ErrCache is returned when there was an error in the backing index. It wraps the original error.
type ErrCache struct {
	Err error
	Msg string
}

// Unwrap returns the underlying error.
func (e ErrCache) Unwrap() error {
	return e.Err
}

// Error returns the message of the current error.
func (e ErrCache) Error() string {
	return e.Msg
}
