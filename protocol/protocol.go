package protocol

import (
	"context"

	t "github.com/ipfs-search/ipfs-search/types"
)

// Protocol represents the interface with one or multiple protocols. It is concurrency-safe.
type Protocol interface {
	GatewayURL(*t.AnnotatedResource) string
	Stat(context.Context, *t.AnnotatedResource) error
	Ls(context.Context, *t.AnnotatedResource, chan<- t.AnnotatedResource) error
}
