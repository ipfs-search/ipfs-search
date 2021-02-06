package types

// LinkType represents the type of a Link as a string.
type LinkType string

// Values for LinkTypes.
const (
	DirectoryLinkType   LinkType = "Directory"
	FileLinkType        LinkType = "File"
	UnknownLinkType     LinkType = "Unknown"
	UnsupportedLinkType LinkType = "Unsupported"
)

// Link from a Document to other Documents.
type Link struct {
	Hash string   `json:"Hash"`
	Name string   `json:"Name"`
	Size uint64   `json:"Size"`
	Type LinkType `json:"Type"`
}

// Links is a collection of links to other Documents.
type Links []Link

// Directory represents a directory resource in an Index.
type Directory struct {
	Document

	Links Links `json:"links"`
}
