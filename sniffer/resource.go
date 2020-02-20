package sniffer

import "fmt"

// Resource represents a resource on the dweb.
type Resource struct {
	Protocol string
	Id       string
}

// URI returns a unique identifier for the resource.
func (r *Resource) URI() string {
	return fmt.Sprintf("%s://%s", r.Protocol, r.Id)
}

// String defaults to the URI
func (r *Resource) String() string {
	return r.URI()
}
