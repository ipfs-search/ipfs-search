package crawler

import (
	"context"
	"time"

	index_types "github.com/ipfs-search/ipfs-search/index/types"
	t "github.com/ipfs-search/ipfs-search/types"
)

var (
	minUpdateAge = time.Duration(time.Hour)
)

func appendReference(refs index_types.References, r *t.Reference) (index_types.References, bool) {
	if r.Parent == nil {
		// No new reference, not updating
		return refs, false
	}

	for _, indexedRef := range refs {
		if indexedRef.ParentHash == r.Parent.ID && indexedRef.Name == r.Name {
			// Existing reference, not updating
			return refs, false
		}
	}

	return append(refs, index_types.Reference{
		ParentHash: r.Parent.ID,
		Name:       r.Name,
	}), true
}

// updateExisting updates known existing items.
func (c *Crawler) updateExisting(ctx context.Context, i *existingItem) error {
	refs, refUpdated := appendReference(i.References, &i.AnnotatedResource.Reference)

	now := time.Now()
	isRecent := now.Sub(i.LastSeen) > minUpdateAge

	if refUpdated || isRecent {
		return i.Index.Update(ctx, i.AnnotatedResource.ID, &index_types.Update{
			LastSeen:   now,
			References: refs,
		})
	}

	return nil
}

// updateMaybeExisting updates an item when it exists and returnes true when item exists.
func (c *Crawler) updateMaybeExisting(ctx context.Context, r *t.AnnotatedResource) (bool, error) {
	existing, err := c.getExistingItem(ctx, r)
	if err != nil {
		return false, err
	}

	// Process existing item
	if existing != nil {
		if existing.Index == c.indexes.Invalids {
			// Already indexed as invalid; we're done
			return true, nil
		}

		// Update item and we're done.
		if err = c.updateExisting(ctx, existing); err != nil {
			return true, err
		}

		return true, nil
	}

	return false, nil
}
