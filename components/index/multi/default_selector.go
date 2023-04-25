package multi

var DefaultSelector = &PropertySelector{
	PropertyNameParts: []string{"metadata", "Content-Type"},
	Matchers:          DefaultMatchers,
	DefaultIndex:      "ipfs_other",
}
