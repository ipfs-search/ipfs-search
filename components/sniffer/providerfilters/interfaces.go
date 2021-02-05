package providerfilters

import (
	t "github.com/ipfs-search/ipfs-search/types"
)

// Filter takes a Provider and returns true when it is to be included, false
// when not and an error when unexpected condition occur.
type Filter interface {
	Filter(t.Provider) (bool, error)
}
