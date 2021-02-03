package crawler

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ipfs-search/ipfs-search/index"
	indexTypes "github.com/ipfs-search/ipfs-search/index/types"
	t "github.com/ipfs-search/ipfs-search/types"
)

func makeDocument(r *t.AnnotatedResource) indexTypes.Document {
	// TODO: Get this through AnnotatedResource from the sniffer.
	now := time.Now().UTC()

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
		// TODO: Add type as constant field for file/directory and as string for invalids.
		// That natively allows us to know about unsupported types (so that we may index them later when supported).
		// Ref: https://www.elastic.co/guide/en/elasticsearch/reference/master/keyword.html#constant-keyword-field-type
		// Type: r.Type,
	}
}

func (c *Crawler) indexInvalid(ctx context.Context, r *t.AnnotatedResource, err error) error {
	// Index unsupported items as invalid.
	return c.indexes.Invalids.Index(ctx, r.ID, &indexTypes.Invalid{
		Error: err.Error(),
	})
}

func (c *Crawler) index(ctx context.Context, r *t.AnnotatedResource) error {
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
		// TODO: Ensure test coverage.
		err = t.ErrUnsupportedType

	case t.PartialType:
		// Not indexing partials, we're done.
		// TODO: Consider indexing partials to avoid future crawling.
		return nil

	case t.UndefinedType:
		panic("undefined type after Stat call")
	}

	if err != nil {
		if errors.Is(err, t.ErrInvalidResource) {
			log.Printf("Indexing invalid '%v', err: %v", r, err)
			// TODO: Ensure test coverage.
			return c.indexInvalid(ctx, r, err)
		}
		return err
	}

	// Index the result
	return index.Index(ctx, r.ID, properties)
}
