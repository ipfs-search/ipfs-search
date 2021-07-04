package crawler

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/label"

	index_types "github.com/ipfs-search/ipfs-search/components/index/types"
	t "github.com/ipfs-search/ipfs-search/types"
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
	ctx, span := c.Tracer.Start(ctx, "crawler.updateExisting")
	defer span.End()

	refs, refsUpdated := appendReference(i.References, &i.AnnotatedResource.Reference)

	now := time.Now()

	isRecent := now.Sub(i.LastSeen) > c.config.MinUpdateAge

	if refsUpdated || isRecent {
		if span.IsRecording() {
			var reason string

			if refsUpdated {
				reason = "reference-added"
			}

			if isRecent {
				reason = "is-recent"
			}
			span.AddEvent(ctx, "Updating",
				label.String("reason", reason),
				label.Any("new-reference", i.AnnotatedResource.Reference),
				label.Stringer("last-seen", i.LastSeen),
			)
		}

		return i.Index.Update(ctx, i.AnnotatedResource.ID, &index_types.Update{
			LastSeen:   now,
			References: refs,
		})
	}

	span.AddEvent(ctx, "Not updating")

	return nil
}

// deletePartial deletes partial items.
func (c *Crawler) deletePartial(ctx context.Context, i *existingItem) error {
	return c.indexes.Partials.Delete(ctx, i.AnnotatedResource.ID)
}

// updateMaybeExisting updates an item when it exists and returnes true when item exists.
func (c *Crawler) updateMaybeExisting(ctx context.Context, r *t.AnnotatedResource) (bool, error) {
	ctx, span := c.Tracer.Start(ctx, "crawler.updateMaybeExisting")
	defer span.End()

	existing, err := c.getExistingItem(ctx, r)
	if err != nil {
		return false, err
	}

	// Process existing item
	if existing != nil {
		if span.IsRecording() {
			span.AddEvent(ctx, "existing", label.Any("index", existing.Index))
		}

		switch existing.Index {
		case c.indexes.Invalids:
			// Already indexed as invalid; we're done
			return true, nil
		case c.indexes.Partials:
			// Found in partials index; previously recognized as an unreferenced partial

			if r.Reference.Parent == nil {
				// Skip unreferenced partial
				return true, nil
			}

			// Referenced partial; delete as partial
			if err = c.deletePartial(ctx, existing); err != nil {
				return true, err
			}

			// Index item as new
			return false, err
		}

		// Update item and we're done.
		if err = c.updateExisting(ctx, existing); err != nil {
			return true, err
		}

		return true, nil
	}

	return false, nil
}
