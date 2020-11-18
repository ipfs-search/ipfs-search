package extractor

import (
	"context"

	t "github.com/ipfs-search/ipfs-search/types"
)

// Extract metadata from a (potentially) referenced resource, updating
// Metadata or returning an error.
type Extractor interface {
	Extract(context.Context, *t.ReferencedResource, t.Metadata) error
}
