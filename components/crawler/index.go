package crawler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/ipfs-search/ipfs-search/components/extractor"
	"github.com/ipfs-search/ipfs-search/components/index"
	indexTypes "github.com/ipfs-search/ipfs-search/components/index/types"
	t "github.com/ipfs-search/ipfs-search/types"
)

func makeDocument(r *t.AnnotatedResource) indexTypes.Document {
	now := time.Now().UTC()

	// Strip milliseconds to cater to legacy ES index format.
	// This can be safely removed after the next reindex with _nomillis removed from time format.
	now = now.Truncate(time.Second)

	var references []indexTypes.Reference
	if r.Reference.Parent != nil {
		references = []indexTypes.Reference{
			{
				ParentHash: r.Reference.Parent.ID,
				Name:       r.Reference.Name,
			},
		}
	}

	// Common Document properties
	return indexTypes.Document{
		FirstSeen:  now,
		LastSeen:   now,
		References: references,
		Size:       r.Size,
	}
}

func (c *Crawler) indexInvalid(ctx context.Context, r *t.AnnotatedResource, err error) error {
	// Index unsupported items as invalid.
	return c.indexes.Invalids.Index(ctx, r.ID, &indexTypes.Invalid{
		Error: err.Error(),
	})
}

func (c *Crawler) getFileProperties(ctx context.Context, r *t.AnnotatedResource) (interface{}, error) {
	var err error

	span := trace.SpanFromContext(ctx)

	properties := &indexTypes.File{
		Document: makeDocument(r),
	}

	// Note; this assumes to be sequential to allow for dependencies amongst extractors.
	for _, e := range c.extractors {
		err = e.Extract(ctx, r, properties)
		if errors.Is(err, extractor.ErrFileTooLarge) {
			// Interpret files which are too large as invalid resources; prevent repeated attempts.
			span.RecordError(err)
			return nil, fmt.Errorf("%w: %v", t.ErrInvalidResource, err)
		}
	}

	return properties, err
}

func (c *Crawler) getDirectoryProperties(ctx context.Context, r *t.AnnotatedResource) (interface{}, error) {
	properties := &indexTypes.Directory{
		Document: makeDocument(r),
	}
	err := c.crawlDir(ctx, r, properties)

	return properties, err
}

func (c *Crawler) getProperties(ctx context.Context, r *t.AnnotatedResource) (index.Index, interface{}, error) {
	var err error

	span := trace.SpanFromContext(ctx)

	switch r.Type {
	case t.FileType:
		f, err := c.getFileProperties(ctx, r)

		return c.indexes.Files, f, err

	case t.DirectoryType:
		d, err := c.getDirectoryProperties(ctx, r)

		return c.indexes.Directories, d, err

	case t.UnsupportedType:
		// Index unsupported items as invalid.
		err = t.ErrUnsupportedType
		span.RecordError(err)

		return nil, nil, err

	case t.PartialType:
		// Index partial (no properties)
		return c.indexes.Partials, &indexTypes.Partial{}, nil

	case t.UndefinedType:
		panic("undefined type after Stat call")

	default:
		panic("unexpected type")
	}
}

func (c *Crawler) index(ctx context.Context, r *t.AnnotatedResource) error {
	ctx, span := c.Tracer.Start(ctx, "crawler.index",
		trace.WithAttributes(attribute.Stringer("type", r.Type)),
	)
	defer span.End()

	index, properties, err := c.getProperties(ctx, r)

	if err != nil {
		if errors.Is(err, t.ErrInvalidResource) {
			log.Printf("Indexing invalid '%v', err: %v", r, err)
			span.RecordError(err)
			return c.indexInvalid(ctx, r, err)
		}

		return err
	}

	// Index the result
	return index.Index(ctx, r.ID, properties)
}
