// Package crawler is grouped around the Crawler component, crawling and indexing content from an AnnotatedResource.
package crawler

import (
	"context"
	"errors"
	"log"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"

	"github.com/ipfs-search/ipfs-search/components/extractor"
	"github.com/ipfs-search/ipfs-search/components/protocol"

	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
)

// Crawler allows crawling of resources.
type Crawler struct {
	config     *Config
	indexes    *Indexes
	queues     *Queues
	protocol   protocol.Protocol
	extractors []extractor.Extractor

	*instr.Instrumentation
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
	ctx, span := c.Tracer.Start(ctx, "crawler.Crawl",
		trace.WithAttributes(label.String("cid", r.ID)),
	)
	defer span.End()

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
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		return err
	}

	if exists {
		log.Printf("Not updating existing resource %v", r)
		span.AddEvent(ctx, "Not updating existing resource")
		return nil
	}

	if err := c.ensureType(ctx, r); err != nil {
		if errors.Is(err, t.ErrInvalidResource) {
			// Resource is invalid, index as such, throwing away ErrInvalidResource in favor of the result of indexing operation.
			log.Printf("Indexing invalid resource %v", r)
			span.AddEvent(ctx, "Indexing invalid resource")

			err = c.indexInvalid(ctx, r, err)
		}

		// Errors from ensureType imply that no type could be found, hence we can't index.
		if err != nil {
			span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		}
		return err
	}

	log.Printf("Indexing new item %v", r)
	err = c.index(ctx, r)
	if err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
	}
	return err
}

// New instantiates a Crawler.
func New(config *Config, indexes *Indexes, queues *Queues, protocol protocol.Protocol, extractors []extractor.Extractor, i *instr.Instrumentation) *Crawler {
	return &Crawler{
		config,
		indexes,
		queues,
		protocol,
		extractors,
		i,
	}
}

func (c *Crawler) ensureType(ctx context.Context, r *t.AnnotatedResource) error {
	ctx, span := c.Tracer.Start(ctx, "crawler.ensureType")
	defer span.End()

	var err error

	if r.Type == t.UndefinedType {
		ctx, cancel := context.WithTimeout(ctx, c.config.StatTimeout)
		defer cancel()

		err = c.protocol.Stat(ctx, r)
		if err != nil {
			span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		}
	}

	return err
}
