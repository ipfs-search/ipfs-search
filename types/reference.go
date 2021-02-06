package types

// Reference to indexed item
type Reference struct {
	Parent *Resource
	Name   string
}

// String shows the name
func (r *Reference) String() string {
	return r.Name
}
