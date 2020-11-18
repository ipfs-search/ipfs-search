package types

type Links struct {
	Hash string `json:"Hash"`
	Name string `json:"Name"`
	Size int    `json:"Size"`
	Type string `json:"Type"`
}

type Directory struct {
	Document

	Links []Links `json:"links"`
}
