package types

import (
	"fmt"
)

// Resource represents a resource on the dweb.
type Resource struct {
	Protocol        // Protocol, e.g. IPFSProtocol
	ID       string // Resource identifier (e.g. CID) for particular Protocol.
}

// URI returns a unique identifier for the resource.
// TODO: Move to dweb-based addresses, once standardized.
func (r *Resource) URI() string {
	return fmt.Sprintf("%s://%s", r.Protocol, r.ID)
}

// String defaults to the URI
func (r *Resource) String() string {
	return r.URI()
}

// IsValid returns true when resource contains a valid value.
func (r *Resource) IsValid() bool {
	return r.Protocol != InvalidProtocol && r.ID != ""
}
