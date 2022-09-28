package bulkgetter

import (
	"context"
	"fmt"
)

// GetRequest represents an item to GET.
type GetRequest struct {
	Index      string
	DocumentID string
	Fields     []string
}

func (r *GetRequest) String() string {
	return fmt.Sprintf("index: %s, id: %s", r.Index, r.DocumentID)
}

// GetResponse represents the response from a GetRequest.
type GetResponse struct {
	Found bool
	Error error
}

// AsyncGetter is an interface to allow for asynchronous getting.
type AsyncGetter interface {
	Get(context.Context, *GetRequest, interface{}) <-chan GetResponse
	Work(context.Context) error
}
