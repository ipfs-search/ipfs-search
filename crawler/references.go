package crawler

import (
	"github.com/ipfs-search/ipfs-search/indexer"
	"log"
)

// updateReferences updates references with name, parentHash and parentName. Returns true when updated
func (i *Indexable) updateReferences(references []indexer.Reference) ([]indexer.Reference, bool) {
	if references == nil {
		// Initialize empty references when none have been found
		references = []indexer.Reference{}
	}

	if i.ParentHash == "" {
		// No parent hash for item, not adding reference
		return references, false
	}

	for _, reference := range references {
		if reference.ParentHash == i.ParentHash {
			// Reference exists, not updating
			return references, false
		}
	}

	// New references found, updating references
	references = append(references, indexer.Reference{
		Name:       i.Name,
		ParentHash: i.ParentHash,
	})

	return references, true
}

// indexReferences retreives or creates references for this hashable,
// returning the resulting references and whether or not the item was
// previously present in the index.
func (i *Indexable) indexReferences() ([]indexer.Reference, bool, error) {
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

	references, referencesUpdated := i.updateReferences(references)

	if alreadyIndexed {
		if referencesUpdated {
			log.Printf("Found %s, reference added: '%s' from %s", i.Hash, i.Name, i.ParentHash)

			properties := metadata{
				"references": references,
			}

			err := i.Indexer.IndexItem(itemType, i.Hash, properties)
			if err != nil {
				return nil, false, err
			}
		} else {
			log.Printf("Found %s, references not updated.", i.Hash)
		}
	} else if referencesUpdated {
		log.Printf("Adding %s, reference '%s' from %s", i.Hash, i.Name, i.ParentHash)
	}

	return references, alreadyIndexed, nil
}
