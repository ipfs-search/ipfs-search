package types

import (
	"fmt"
)

// AnnotatedResource annotates a referenced Resource with additional information.
type AnnotatedResource struct {
	*Resource
	Reference `json:",omitempty"`
	Stat      `json:",omitempty"`
}

// String returns the first reference or the URI.
func (r *AnnotatedResource) String() string {
	if r.Reference.Name != "" {
		return fmt.Sprintf("%s (%s)", r.Reference.Name, r.URI())
	}

	return r.URI()
}
