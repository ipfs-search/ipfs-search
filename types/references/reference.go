package references

// Reference to indexed item
type Reference struct {
	ParentHash string `json:"parent_hash"`
	Name       string `json:"name"`
}

// String shows the name
func (r *Reference) String() string {
	return r.Name
}
