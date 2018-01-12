package crawler

import (
	"fmt"
	"github.com/ipfs-search/ipfs-search/indexer"
)

// hashURL returns the IPFS URL for a particular hash
func hashURL(hash string) string {
	return fmt.Sprintf("/ipfs/%s", hash)
}

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
