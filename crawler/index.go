package crawler

import (
	"context"
	"time"

	"github.com/ipfs-search/ipfs-search/index"
	index_types "github.com/ipfs-search/ipfs-search/index/types"
	t "github.com/ipfs-search/ipfs-search/types"
)

func (c *Crawler) index(ctx context.Context, r *t.AnnotatedResource) error {
	var err error
	var index index.Index
	var properties interface{}

	// TODO: Get this through AnnotatedResource from the sniffer.
	now := time.Now().UTC()

	// Common Document properties
	document := index_types.Document{
		FirstSeen: now,
		LastSeen:  now,
		References: []index_types.Reference{
			{
				ParentHash: r.Reference.Parent.ID,
				Name:       r.Reference.Name,
			},
		},
		Size: r.Size,
		// TODO: Add this as constant field for file/directory and as string for invalids.
		// That natively allows us to know about unsupported types (so that we may index them later when supported).
		// Ref: https://www.elastic.co/guide/en/elasticsearch/reference/master/keyword.html#constant-keyword-field-type
		// Type: r.Type,
	}

	switch r.Type {
	case t.UndefinedType:
		panic("undefined type after Stat call")

	case t.UnsupportedType:
		// (For now) index this as invalid
		index = c.indexes.Invalid

	case t.FileType:
		f := index_types.File{
			Document: document,
		}
		err = c.extractor.Extract(ctx, r, &f)
		// TODO
		// This might yield an unexpected end of file, or a non-timeout or IPFS daemon related
		// error which renders a file invalid!

		index = c.indexes.Files
		properties = f

	case t.DirectoryType:
		d := index_types.Directory{
			Document: document,
		}
		err = c.crawlDirectory(ctx, r, &d)
		// TODO
		// Depending on err, this might be invalid!
		// (although unlikely, as stat was above)

		index = c.indexes.Directories
		properties = d
	}
	if err != nil {
		return err
	}

	// Index the result
	return index.Index(ctx, r.Resource.ID, properties)
}
