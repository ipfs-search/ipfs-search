package types

// Invalid represents invalid (unindexable) resources in an Index.
type Invalid struct {
	Error string `json:"error"`
}
