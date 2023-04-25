package multi

var DefaultSelector = &PropertySelector{
	PropertyName: "metadata.Content-Type",
	Matchers:     DefaultMatchers,
	DefaultIndex: "ipfs_other",
}
