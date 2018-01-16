package crawler

import (
	"github.com/ipfs-search/ipfs-search/indexer"
	"log"
)

// updateReferences updates references with name, parentHash and parentName. Returns true when updated
func updateReferences(references []indexer.Reference, name string, parentHash string) ([]indexer.Reference, bool) {
	if references == nil {
		// Initialize empty references when none have been found
		references = []indexer.Reference{}
	}

	if parentHash == "" {
		// No parent hash for item, not adding reference
		return references, false
	}

	for _, reference := range references {
		if reference.ParentHash == parentHash {
			// Reference exists, not updating
			return references, false
		}
	}

	references = append(references, indexer.Reference{
		Name:       name,
		ParentHash: parentHash,
	})

	return references, true
}

func (c *Crawler) indexReferences(hash string, name string, parentHash string) ([]indexer.Reference, bool, error) {
	var alreadyIndexed bool

	references, itemType, err := c.Indexer.GetReferences(hash)
	if err != nil {
		return nil, false, err
	}

	// TODO: Handle this more explicitly, use and detect NotFound
	if references == nil {
		alreadyIndexed = false
	} else {
		alreadyIndexed = true
	}

	references, referencesUpdated := updateReferences(references, name, parentHash)

	if alreadyIndexed {
		if referencesUpdated {
			log.Printf("Found %s, reference added: '%s' from %s", hash, name, parentHash)

			properties := map[string]interface{}{
				"references": references,
			}

			err := c.Indexer.IndexItem(itemType, hash, properties)
			if err != nil {
				return nil, false, err
			}
		} else {
			log.Printf("Found %s, references not updated.", hash)
		}
	} else if referencesUpdated {
		log.Printf("Adding %s, reference '%s' from %s", hash, name, parentHash)
	}

	return references, alreadyIndexed, nil
}
