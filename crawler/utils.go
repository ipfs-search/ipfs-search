package crawler

import (
	"fmt"
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
