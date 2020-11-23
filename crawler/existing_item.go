package crawler

import (
	"context"
	"time"

	"github.com/ipfs-search/ipfs-search/index"
	index_types "github.com/ipfs-search/ipfs-search/index/types"
	t "github.com/ipfs-search/ipfs-search/types"
)

type ExistingItem struct {
	*t.AnnotatedResource
	index.Index
	index_types.References
	LastSeen time.Time
}

func (c *Crawler) getExistingItem(context.Context, *t.AnnotatedResource) (*ExistingItem, error) {
	// TODO

	// refs := new(index_types.References)
	// exists, err := i.Get(ctx, r.ID, refs, "references")

	// if err != nil {
	// 	return false, err
	// }

	// if exists {

	// 	// Reference already existing, not updating
	// }

	// // Not existing, not updating
	// return false, nil

	return nil, nil
}
