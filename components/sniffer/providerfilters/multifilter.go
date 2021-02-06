package providerfilters

import (
	"errors"
	"fmt"

	t "github.com/ipfs-search/ipfs-search/types"
)

var (
	errFilter = errors.New("filter error")
)

// MultiFilter efficiently combines multiple filters into a single filter.
type MultiFilter struct {
	filters []Filter
}

// NewMultiFilter returns a pointer to a new MultiFilter.
func NewMultiFilter(filters ...Filter) *MultiFilter {
	return &MultiFilter{
		filters,
	}
}

// Filter takes a Provider and returns true when it is to be included, false
// when not and an error when unexpected condition occur.
func (m *MultiFilter) Filter(p t.Provider) (bool, error) {
	for _, f := range m.filters {
		include, err := f.Filter(p)

		if err != nil {
			return false, fmt.Errorf("%w: %s with filter %v", errFilter, err, f)
		}

		if !include {
			return false, nil
		}
	}

	return true, nil
}
