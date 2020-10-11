package types

import (
	"fmt"
	"go.opentelemetry.io/otel/api/trace"
)

// Resource represents a resource on the dweb.
type Resource struct {
	Protocol    string            // Protocol as a string, e.g. "ipfs"
	ID          string            // Identifier for this resource, unique together with Protocol
	SpanContext trace.SpanContext // SpanContext allows a Resource' processing to be traceable across the program
}

// URI returns a unique identifier for the resource.
func (r *Resource) URI() string {
	return fmt.Sprintf("%s://%s", r.Protocol, r.ID)
}

// String defaults to the URI
func (r *Resource) String() string {
	return r.URI()
}
