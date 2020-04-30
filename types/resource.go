package types

import "fmt"

// Resource represents a resource on the dweb.
type Resource struct {
	Protocol string // Protocol as a string, e.g. "ipfs"
	Id       string // Identifier for this resource, unique together with Protocol
}

// URI returns a unique identifier for the resource.
func (r *Resource) URI() string {
	return fmt.Sprintf("%s://%s", r.Protocol, r.Id)
}

// String defaults to the URI
func (r *Resource) String() string {
	return r.URI()
}
