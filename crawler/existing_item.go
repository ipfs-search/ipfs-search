package crawler

import (
	"github.com/ipfs-search/ipfs-search/indexer"
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
func (i *existingItem) updateIndex() error {
	properties := metadata{
		"references": i.references,
		"last-seen":  nowISO(),
	}

	return i.Indexer.IndexItem(i.itemType, i.Hash, properties)
}

// update updates existing items (if they in fact do exist)
func (i *existingItem) update() error {
	i.updateReferences()

	if i.exists {
		return i.updateIndex()
	}

	return nil
}

// getExistingItem returns existingItem from index
func (i *Indexable) getExistingItem() (*existingItem, error) {
	references, itemType, err := i.Indexer.GetReferences(i.Hash)
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
	return !i.skipItem() || !i.exists

}
