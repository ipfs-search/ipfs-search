package crawler

import (
	"fmt"

	t "github.com/ipfs-search/ipfs-search/types"
)

// UnexpectedTypeError is returned when encountering an unexpected resource type.
type UnexpectedTypeError struct {
	t.ResourceType
}

func (e UnexpectedTypeError) Error() string {
	return fmt.Sprintf("unexpected type: %s", e.ResourceType)
}
