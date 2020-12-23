package types

type AnnotatedResource struct {
	*Resource
	Reference `json:",omitempty"`
	Stat      `json:",omitempty"`
}
