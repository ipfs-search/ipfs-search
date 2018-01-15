package crawler

import (
	"fmt"
	"github.com/ipfs-search/ipfs-search/indexer"
)

// hashURL returns the IPFS URL for a particular hash
func hashURL(hash string) string {
	return fmt.Sprintf("/ipfs/%s", hash)
}

// filenameURL returns an IPFS reference including a filename, if available.
// e.g. /ipfs/<parent_hash>/my_file.jpg instead of /ipfs/<file_hash>/
// This helps Tika with file type detection.
func filenameURL(args *Args) (path string) {
	if args.Name != "" && args.ParentHash != "" {
		return fmt.Sprintf("/ipfs/%s/%s", args.ParentHash, args.Name)
	}

	// No name & parent hash available
	return fmt.Sprintf("/ipfs/%s", args.Hash)
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
