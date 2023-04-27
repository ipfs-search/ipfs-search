package worker

import (
	"context"

	t "github.com/ipfs-search/ipfs-search/types"
)

// Crawler represents the public interface of a crawler.
// TODO: Integrate/refactor crawler component.
type Crawler interface {
	Crawl(ctx context.Context, r *t.AnnotatedResource) error
}
