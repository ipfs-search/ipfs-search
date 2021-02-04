package types

// AnnotatedResource annotates a referenced Resource with additional information.
type AnnotatedResource struct {
	*Resource
	Reference `json:",omitempty"`
	Stat      `json:",omitempty"`
}
