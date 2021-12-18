package crawler

import (
	"context"
	"fmt"
	"log"
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

	switch i.Source {
	case t.DirectorySource:
		// Item referenced from a directory, consider updating references (but not last-seen).
		refs, refsUpdated := appendReference(i.References, &i.AnnotatedResource.Reference)

		if refsUpdated {
			span.AddEvent(ctx, "Updating",
				label.String("reason", "reference-added"),
				label.Any("new-reference", i.AnnotatedResource.Reference),
			)

			return i.Index.Update(ctx, i.AnnotatedResource.ID, &index_types.Update{
				References: refs,
			})
		}

	case t.SnifferSource, t.UnknownSource:
		// TODO: Remove UnknownSource after sniffer is updated and queue is flushed.
		// Item sniffed, conditionally update last-seen.
		now := time.Now()
		isRecent := now.Sub(*i.LastSeen) > c.config.MinUpdateAge

		if isRecent {
			span.AddEvent(ctx, "Updating",
				label.String("reason", "is-recent"),
				label.Stringer("last-seen", i.LastSeen),
			)

			return i.Index.Update(ctx, i.AnnotatedResource.ID, &index_types.Update{
				LastSeen: &now,
			})
		}

	default:
		// Panic for unexpected Source values, instead of hard failing.
		panic(fmt.Sprintf("Unexpected source %s for item %+v", i.Source, i))
	}

	span.AddEvent(ctx, "Not updating")

	return nil
}

// deletePartial deletes partial items.
func (c *Crawler) deletePartial(ctx context.Context, i *existingItem) error {
	return c.indexes.Partials.Delete(ctx, i.ID)
}

// processPartial processes partials found in index; previously recognized as an unreferenced partial
func (c *Crawler) processPartial(ctx context.Context, i *existingItem) (bool, error) {
	if i.Reference.Parent == nil {
		log.Printf("Quick-skipping unreferenced partial %v", i)

		// Skip unreferenced partial
		return true, nil
	}

	// Referenced partial; delete as partial
	if err := c.deletePartial(ctx, i); err != nil {
		return true, err
	}

	// Index item as new
	return false, nil
}

func (c *Crawler) processExisting(ctx context.Context, i *existingItem) (bool, error) {
	switch i.Index {
	case c.indexes.Invalids:
		// Already indexed as invalid; we're done
		return true, nil
	case c.indexes.Partials:
		return c.processPartial(ctx, i)
	}

	// Update item and we're done.
	if err := c.updateExisting(ctx, i); err != nil {
		return true, err
	}

	return true, nil
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

		return c.processExisting(ctx, existing)
	}

	return false, nil
}
