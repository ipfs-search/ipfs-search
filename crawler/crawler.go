package crawler

import (
	"context"

	"github.com/ipfs-search/ipfs-search/extractor"
	"github.com/ipfs-search/ipfs-search/protocol"
	t "github.com/ipfs-search/ipfs-search/types"
)

// Crawler allows crawling of resources.
type Crawler struct {
	config    *Config
	indexes   Indexes
	queues    Queues
	protocol  protocol.Protocol
	extractor extractor.Extractor
}

// Crawl updates existing or crawls new resources, extracting metadata where applicable.
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

	if err := c.ensureType(ctx, r); err != nil {
		if c.protocol.IsInvalidResourceErr(err) {
			// Resource is invalid, index as such, overwriting previous error.
			err = c.indexInvalid(ctx, r, err)
		}

		// Errors from ensureType imply that no type could be found, hence we can't index.
		return err
	}

	// Index new item
	return c.index(ctx, r)
}

// New instantiates a Crawler.
func New(config *Config, indexes Indexes, queues Queues, protocol protocol.Protocol, extractor extractor.Extractor) *Crawler {
	return &Crawler{
		config,
		indexes,
		queues,
		protocol,
		extractor,
	}
}

func (c *Crawler) ensureType(ctx context.Context, r *t.AnnotatedResource) error {
	if r.Type == t.UndefinedType {
		ctx, cancel := context.WithTimeout(ctx, c.config.StatTimeout)
		defer cancel()

		return c.protocol.Stat(ctx, r)
	}

	return nil
}
