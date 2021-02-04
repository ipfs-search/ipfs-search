package crawler

import (
	"context"
	"errors"
	"log"

	"github.com/ipfs-search/ipfs-search/extractor"
	"github.com/ipfs-search/ipfs-search/protocol"
	t "github.com/ipfs-search/ipfs-search/types"
)

// Crawler allows crawling of resources.
type Crawler struct {
	config    *Config
	indexes   *Indexes
	queues    *Queues
	protocol  protocol.Protocol
	extractor extractor.Extractor
}

func isSupportedType(rType t.ResourceType) bool {
	switch rType {
	case t.UndefinedType, t.FileType, t.DirectoryType:
		return true
	default:
		return false
	}
}

// Crawl updates existing or crawls new resources, extracting metadata where applicable.
func (c *Crawler) Crawl(ctx context.Context, r *t.AnnotatedResource) error {
	var err error

	if r.Protocol == t.InvalidProtocol {
		// Sending items with an invalid protocol to Crawl() is a programming error and
		// should never happen.
		panic("invalid protocol")
	}

	if !isSupportedType(r.Type) {
		// Calling crawler with unsupported types is undefined behaviour.
		panic("invalid type for crawler")
	}

	exists, err := c.updateMaybeExisting(ctx, r)
	if err != nil {
		return err
	}

	if exists {
		log.Printf("Not updating existing resource %v", r)
		return nil
	}

	if err := c.ensureType(ctx, r); err != nil {
		if errors.Is(err, t.ErrInvalidResource) {
			// Resource is invalid, index as such, throwing away ErrInvalidResource in favor of the result of indexing operation.
			log.Printf("Indexing invalid resource %v", r)

			err = c.indexInvalid(ctx, r, err)
		}

		// Errors from ensureType imply that no type could be found, hence we can't index.
		return err
	}

	log.Printf("Indexing new item %v", r)
	return c.index(ctx, r)
}

// New instantiates a Crawler.
func New(config *Config, indexes *Indexes, queues *Queues, protocol protocol.Protocol, extractor extractor.Extractor) *Crawler {
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
