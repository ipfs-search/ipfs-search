package filters

import (
	t "github.com/ipfs-search/ipfs-search/types"
)

type multiFilter struct {
	filters []Filter
}

// MultiFilter efficiently combines multiple filters into a single one
func MultiFilter(filters ...Filter) *multiFilter {
	return &multiFilter{
		filters,
	}
}

// Filter returns false for the first filter returning false, true otherwise
func (m *multiFilter) Filter(p t.Provider) (bool, error) {
	for _, f := range m.filters {
		include, err := f.Filter(p)

		if err != nil {
			return false, t.NewProviderErrorf(err, p, "Error %s with filter %v", err, f)
		}

		if !include {
			return false, nil
		}
	}

	return true, nil
}
