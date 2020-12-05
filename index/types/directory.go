package types

type LinkType string

const (
	DirectoryLinkType = "Directory"
	FileLinkType      = "File"
)

type Link struct {
	Hash string   `json:"Hash"`
	Name string   `json:"Name"`
	Size uint64   `json:"Size"`
	Type LinkType `json:"Type"`
}

type Links []Link

type Directory struct {
	Document

	Links Links `json:"links"`
}
