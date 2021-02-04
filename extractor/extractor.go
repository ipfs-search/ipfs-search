// Package extractor is grouped around the Extractor component, extracting metadata from an AnnotatedResource.
package extractor

import (
	"context"

	t "github.com/ipfs-search/ipfs-search/types"
)

// Extractor extracts metadata from an AnnotatedResource, updating the metadata interface or returning an error.
type Extractor interface {
	Extract(ctx context.Context, resource *t.AnnotatedResource, metadata interface{}) error
}
