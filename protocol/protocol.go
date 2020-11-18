package protocol

import (
	"context"

	t "github.com/ipfs-search/ipfs-search/types"
)

// Protocol represents the interface with one or multiple protocols. It is concurrency-safe.
type Protocol interface {
	GatewayURL(*t.ReferencedResource) string
	Stat(context.Context, *t.Resource) (*t.ReferencedResource, error)
	Ls(context.Context, *t.Resource, chan<- t.ReferencedResource) error
}
