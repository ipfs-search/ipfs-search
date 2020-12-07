package crawler

import (
	"context"

	"github.com/ipfs-search/ipfs-search/extractor"
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

	if r.Protocol == t.InvalidProtocol {
		// Sending items with an invalid protocol to Crawl() is a programming error and
		// should never happen.
		panic("invalid protocol")
	}

	switch r.Type {
	case t.UnsupportedType, t.PartialType:
		// Crawler should never be called with these types, this is unsupported behaviour.
		panic("invalid type for crawler")
	}

	exists, err := c.updateMaybeExisting(ctx, r)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	// Ensure type is present
	if r.Type == t.UndefinedType {
		// Get size and type

		// TODO: Implement a timeout for Stat call here.

		if err := c.protocol.Stat(ctx, r); err != nil {
			if c.protocol.IsInvalidResourceErr(err) {
				// Resource is invalid, index as such, overwriting previous error.
				return c.indexInvalid(ctx, r, err)
			}
			return err
		}
	}

	// Index new item
	return c.index(ctx, r)
}

func New(indexes Indexes, queues Queues, protocol protocol.Protocol, extractor extractor.Extractor) *Crawler {
	return &Crawler{
		indexes,
		queues,
		protocol,
		extractor,
	}
}
