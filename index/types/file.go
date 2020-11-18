package types

type Language struct {
	Confidence string  `json:"confidence"`
	Language   string  `json:"language"`
	RawScore   float64 `json:"rawScore"`
}

type Metadata map[string]interface{}

type File struct {
	Document

	Content         string   `json:"content"`
	IpfsTikaVersion string   `json:"ipfs_tika_version"`
	Language        Language `json:"language"`
	Metadata        Metadata `json:"metadata"`
	Urls            []string `json:"urls"`
}
