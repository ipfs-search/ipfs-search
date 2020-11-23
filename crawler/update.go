package crawler

import (
	"context"
	"fmt"

	index_types "github.com/ipfs-search/ipfs-search/index/types"
	t "github.com/ipfs-search/ipfs-search/types"
)

func appendReference(refs index_types.References, r *t.Reference) (index_types.References, bool) {
	if r.Parent == nil {
		// No new reference, not updating
		return refs, false
	}

	for _, indexedRef := range refs {
		if r.Parent.Protocol != t.IPFSProtocol {
			panic(fmt.Sprintf("Unsupported protocol: %s", r.Parent.Protocol))
		}

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

func (c *Crawler) update(ctx context.Context, i *ExistingItem) error {
	refs, updated := appendReference(i.References, &i.AnnotatedResource.Reference)
	if updated {
		// Updated references, updating in index
		return i.Index.Update(ctx, i.AnnotatedResource.ID, refs)
	}

	return nil
}
