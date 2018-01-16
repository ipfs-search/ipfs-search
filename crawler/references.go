package crawler

import (
	"github.com/ipfs-search/ipfs-search/indexer"
	"log"
)

// updateReferences updates references with Name and ParentHash
func (i *Indexable) updateReferences(references []indexer.Reference) []indexer.Reference {
	if references == nil {
		// Initialize empty references when none have been found
		references = []indexer.Reference{}
	}

	if i.ParentHash == "" {
		// No parent hash for item, not adding reference
		return references
	}

	for _, reference := range references {
		if reference.ParentHash == i.ParentHash {
			// Reference exists, not updating
			return references
		}
	}

	// New references found, updating references
	references = append(references, indexer.Reference{
		Name:       i.Name,
		ParentHash: i.ParentHash,
	})

	return references
}

// Return an updated list of references (existing plus new) and whether or the
// item was previously indexed.
func (i *Indexable) getReferences() ([]indexer.Reference, bool, error) {
	var alreadyIndexed bool

	references, itemType, err := i.Indexer.GetReferences(i.Hash)
	if err != nil {
		return nil, false, err
	}

	// TODO: Handle this more explicitly, use and detect NotFound
	if references == nil {
		alreadyIndexed = false
	} else {
		alreadyIndexed = true
	}

	references = i.updateReferences(references)

	return references, alreadyIndexed, nil
}
