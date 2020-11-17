package crawler

import (
	"context"
	"log"

	"github.com/ipfs-search/ipfs-search/index"
	t "github.com/ipfs-search/ipfs-search/types"
)

type existingItem struct {
	*Indexable
	exists     bool
	references t.References
	index      index.Index
}

// updateReferences updates references with Name and ParentHash
func (i *existingItem) updateReferences() {
	newRef := ReferenceFromIndexable(i.Indexable)

	if newRef.ParentHash == "" || i.references.Contains(newRef) {
		// Not updating references
		return
	}

	log.Printf("Adding reference '%v' to %v", newRef, i)
	i.references = append(i.references, *newRef)
}

// updateItem updates references and last seen date
func (i *existingItem) updateIndex(ctx context.Context) error {
	properties := t.Metadata{
		"references": i.references,
		"last-seen":  nowISO(),
	}

	return i.index.Update(ctx, i.Hash, properties)
}

// update updates existing items (if they in fact do exist)
func (i *existingItem) update(ctx context.Context) error {
	if !i.skipItem() {
		// Update references always; this also adds existing to them
		// FIXME: I know, this is bad design...
		i.updateReferences()

		if i.exists {
			log.Printf("Updating %v", i)
			return i.updateIndex(ctx)
		}
	}

	return nil
}

// skipItem determines whether a particular item should not be indexed
// This holds particularly to partial content.
func (i *existingItem) skipItem() bool {
	// TODO; this is currently called in update() and shouldCrawl and
	// yields duplicate output. Todo; make this return an error or nil.
	if i.Size == i.Config.PartialSize && i.ParentHash == "" {
		log.Printf("Skipping unreferenced partial content for item %v", i)
		return true
	}

	if i.index == i.InvalidIndex {
		log.Printf("Skipping update of invalid %v", i)
		return true
	}

	return false
}

// getExistingItem returns existingItem from index
func (i *Indexable) getExistingItem(ctx context.Context) (*existingItem, error) {
	if i == nil {
		panic("Indexable should not be nil")
	}

	indexes := []index.Getter{i.InvalidIndex, i.FileIndex, i.DirectoryIndex}

	// Container for query reference fetch results
	src := &struct {
		references t.References
	}{}

	getterIndex, err := index.MultiGet(ctx, indexes, i.Hash, src, "references")

	if err != nil {
		return nil, err
	}

	// Type assert back to index MulitGetter works on getters
	index, ok := getterIndex.(index.Index)

	if getterIndex != nil && !ok {
		panic("Cannot assert Getter back to Index!")
	}

	item := &existingItem{
		Indexable:  i,
		exists:     getterIndex != nil,
		references: src.references,
		index:      index,
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
