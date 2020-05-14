package references

// References represents a list of references
type References []Reference

// Contains returns true of a given reference exists, false when it doesn't
func (references References) Contains(newRef *Reference) bool {
	for _, r := range references {
		if r.ParentHash == newRef.ParentHash {
			return true
		}
	}

	return false
}
