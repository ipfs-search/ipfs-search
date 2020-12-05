package types

const (
	UnsupportedTypeError = "unsupported type"
)

type Invalid struct {
	Error string `json:"error"`
}
