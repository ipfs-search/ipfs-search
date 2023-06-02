package index

// Factory creates and returns a new Index with given name.
type Factory interface {
	NewIndex(name string) Index
}
