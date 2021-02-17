package crawler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"

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

func (c *Crawler) index(ctx context.Context, r *t.AnnotatedResource) error {
	ctx, span := c.Tracer.Start(ctx, "crawler.index",
		trace.WithAttributes(label.Stringer("type", r.Type)),
	)
	defer span.End()

	var (
		err        error
		index      index.Index
		properties interface{}
	)

	switch r.Type {
	case t.FileType:
		f := &indexTypes.File{
			Document: makeDocument(r),
		}
		err = c.extractor.Extract(ctx, r, f)
		if errors.Is(err, extractor.ErrFileTooLarge) {
			// Interpret files which are too large as invalid resources; prevent repeated attempts.
			span.RecordError(ctx, err)
			err = fmt.Errorf("%w: %v", t.ErrInvalidResource, err)
		}

		index = c.indexes.Files
		properties = f

	case t.DirectoryType:
		d := &indexTypes.Directory{
			Document: makeDocument(r),
		}
		err = c.crawlDir(ctx, r, d)

		index = c.indexes.Directories
		properties = d

	case t.UnsupportedType:
		// Index unsupported items as invalid.
		span.RecordError(ctx, err)
		err = t.ErrUnsupportedType

	case t.PartialType:
		// Not indexing partials (for now), we're done.
		span.AddEvent(ctx, "partial")
		return nil

	case t.UndefinedType:
		panic("undefined type after Stat call")

	default:
		panic("unexpected type")
	}

	if err != nil {
		if errors.Is(err, t.ErrInvalidResource) {
			log.Printf("Indexing invalid '%v', err: %v", r, err)
			span.RecordError(ctx, err)
			return c.indexInvalid(ctx, r, err)
		}

		return err
	}

	// Index the result
	return index.Index(ctx, r.ID, properties)
}
