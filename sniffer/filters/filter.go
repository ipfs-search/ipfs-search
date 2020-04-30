package filters

import (
	t "github.com/ipfs-search/ipfs-search/types"
)

// Filter takes a provider, returning true if it is to be included or false when
// it is to be discarded.
type Filter interface {
	Filter(t.Provider) (bool, error)
}
