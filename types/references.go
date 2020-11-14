package types

import (
	"fmt"
)

// Reference to indexed item
type Reference struct {
	ParentHash string `json:"parent_hash"`
	Name       string `json:"name"`
}

// String shows the name
func (r *Reference) String() string {
	return r.Name
}

// References represents a list of references
type References []Reference

// Contains returns true of a given reference exists, false when it doesn't
func (references References) Contains(newRef *Reference) bool {
	for _, r := range references {
		if r.ParentHash == newRef.ParentHash {
			return true
		}
	}

	return false
}

// ReferencedResource is a resource with zero or more references to it.
type ReferencedResource struct {
	*Resource
	References
}

// GatewayPath returns the path for requesting a resource from an IPFS gateway.
// If a reference is available, it is used to generate the filename to facilitate content
// type detection (e.g. /ipfs/<parent_hash>/my_file.jpg instead of /ipfs/<file_hash>/).
func (r ReferencedResource) GatewayPath() string {
	if len(r.References) > 0 && r.References[0].Name != "" {
		// Named reference, use it for generating path
		return fmt.Sprintf("/ipfs/%s/%s", r.References[0].ParentHash, r.References[0].Name)
	}

	return fmt.Sprintf("/ipfs/%s", r.ID)
}
