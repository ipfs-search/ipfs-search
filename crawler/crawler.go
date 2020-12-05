package crawler

import (
	"context"

	"github.com/ipfs-search/ipfs-search/extractor"
	index_types "github.com/ipfs-search/ipfs-search/index/types"
	"github.com/ipfs-search/ipfs-search/protocol"
	t "github.com/ipfs-search/ipfs-search/types"
)

type Crawler struct {
	indexes   Indexes
	queues    Queues
	protocol  protocol.Protocol
	extractor extractor.Extractor
}

func (c *Crawler) Crawl(ctx context.Context, r *t.AnnotatedResource) error {
	var err error

	existing, err := c.getExistingItem(ctx, r)
	if err != nil {
		return err
	}

	// Process existing item
	if existing != nil {
		if existing.Index == c.indexes.Invalid {
			// Already indexed as invalid; we're done
			return nil
		}

		// Update item and we're done.
		return c.update(ctx, existing)
	}

	// Ensure type is present
	if r.Type == t.UndefinedType {
		// Get size and type
		if err := c.protocol.Stat(ctx, r); err != nil {
			// Depending on error, index as invalid
			return err
		}
	}

	if r.Type == t.PartialType {
		// Not indexing partials, we're done.
		return nil
	}

	// Index new item
	return c.index(ctx, r)
}

func (c *Crawler) crawlDirectory(ctx context.Context, r *t.AnnotatedResource, properties *index_types.Directory) error {
	// TODO

	// queue directory entries

	// update entries in properties

	return nil
}

func New(indexes Indexes, queues Queues, protocol protocol.Protocol, extractor extractor.Extractor) *Crawler {
	return &Crawler{
		indexes,
		queues,
		protocol,
		extractor,
	}
}
