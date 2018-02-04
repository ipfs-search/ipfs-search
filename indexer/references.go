package indexer

// Reference to indexed item
type Reference struct {
	ParentHash string `json:"parent_hash"`
	Name       string `json:"name"`
}

// String shows the name
func (r *Reference) String() string {
	return r.Name
}

// References represents a list of references
type References []Reference

// Exists returns true of a given reference exists, false when it doesn't
func (references References) Contains(newRef *Reference) bool {
	for _, r := range references {
		if r.ParentHash == newRef.ParentHash {
			return true
		}
	}

	return false
}
