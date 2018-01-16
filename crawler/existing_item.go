package crawler

import (
	"context"
	"github.com/ipfs-search/ipfs-search/indexer"
	"log"
)

type existingItem struct {
	*Indexable
	exists     bool
	references []indexer.Reference
	itemType   string
}

// updateReferences updates references with Name and ParentHash
func (i *existingItem) updateReferences() {
	if i.references == nil {
		// Initialize empty references when none have been found
		i.references = []indexer.Reference{}
		return
	}

	if i.ParentHash == "" {
		// No parent hash for item, not adding reference
		return
	}

	for _, reference := range i.references {
		if reference.ParentHash == i.ParentHash {
			// Reference exists, not updating
			return
		}
	}

	// New references found, updating references
	i.references = append(i.references, indexer.Reference{
		Name:       i.Name,
		ParentHash: i.ParentHash,
	})
}

// updateItem updates references (and later also last seen date)
func (i *existingItem) updateIndex(ctx context.Context) error {
	properties := metadata{
		"references": i.references,
		"last-seen":  nowISO(),
	}

	return i.Indexer.IndexItem(ctx, i.itemType, i.Hash, properties)
}

// update updates existing items (if they in fact do exist)
func (i *existingItem) update(ctx context.Context) error {
	if i.exists && !i.skipItem() {
		i.updateReferences()

		log.Printf("Updating %s", i)
		return i.updateIndex(ctx)
	}

	return nil
}

// skipItem determines whether a particular item should not be indexed
// This holds particularly to partial content.
func (i *existingItem) skipItem() bool {
	// TODO; this is currently called in update() and shouldCrawl and
	// yields duplicate output. Todo; make this return an error or nil.
	if i.Size == i.Config.PartialSize && i.ParentHash == "" {
		log.Printf("Skipping unreferenced partial content for item %s", i)
		return true
	}

	if i.itemType == "invalid" {
		log.Printf("Skipping update of invalid %s", i)
		return true
	}

	return false
}

// getExistingItem returns existingItem from index
func (i *Indexable) getExistingItem(ctx context.Context) (*existingItem, error) {
	if i == nil {
		panic("Indexable should not be nil")
	}

	references, itemType, err := i.Indexer.GetReferences(ctx, i.Hash)
	if err != nil {
		return nil, err
	}

	item := &existingItem{
		Indexable:  i,
		exists:     references != nil, // references == nil -> doesn't exist
		references: references,
		itemType:   itemType,
	}

	return item, nil
}

// shouldCrawl returns whether or not this item should be crawled
func (i *existingItem) shouldCrawl() bool {
	if i == nil {
		panic("Existingitem should not be nil")
	}

	return !(i.skipItem() || i.exists)
}
