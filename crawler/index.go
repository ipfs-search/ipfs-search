package crawler

import (
	"context"
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
		// TODO: Add this as constant field for file/directory and as string for invalids.
		// That natively allows us to know about unsupported types (so that we may index them later when supported).
		// Ref: https://www.elastic.co/guide/en/elasticsearch/reference/master/keyword.html#constant-keyword-field-type
		// Type: r.Type,
	}
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
		// TODO
		// This might yield an unexpected end of file, or a non-timeout or IPFS daemon related
		// error which renders a file invalid!

		index = c.indexes.Files
		properties = f

	case t.DirectoryType:
		d := &indexTypes.Directory{
			Document: makeDocument(r),
		}
		err = c.crawlDir(ctx, r, d)
		// TODO
		// Depending on err, this might be invalid!
		// (although unlikely, as stat was above)

		index = c.indexes.Directories
		properties = d

	case t.UnsupportedType:
		// Index unsupported items as invalid.
		index = c.indexes.Invalid
		properties = &indexTypes.Invalid{
			Error: indexTypes.UnsupportedTypeError,
		}

	case t.PartialType:
		// Not indexing partials, we're done.
		// TODO: Consider indexing partials to avoid future crawling.
		return nil

	case t.UndefinedType:
		panic("undefined type after Stat call")
	}

	if err != nil {
		return err
	}

	// Index the result
	return index.Index(ctx, r.ID, properties)
}
