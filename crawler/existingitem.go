package crawler

import (
	"context"

	"github.com/ipfs-search/ipfs-search/index"
	index_types "github.com/ipfs-search/ipfs-search/index/types"
	t "github.com/ipfs-search/ipfs-search/types"
)

type existingItem struct {
	*t.AnnotatedResource
	index.Index
	*index_types.Update
}

func (c *Crawler) getExistingItem(ctx context.Context, r *t.AnnotatedResource) (*existingItem, error) {
	indexes := []index.Index{c.indexes.Files, c.indexes.Directories, c.indexes.Invalids}

	update := new(index_types.Update)

	index, err := index.MultiGet(ctx, indexes, r.ID, update, "references", "last-seen")
	if err != nil {
		return nil, err
	}

	if index == nil {
		// Not found
		return nil, nil
	}

	return &existingItem{
		r, index, update,
	}, nil
}
