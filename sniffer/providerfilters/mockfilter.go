package providerfilters

import (
	"github.com/ipfs-search/ipfs-search/types"
)

// MockFilter represents a mock for a Filter.
type MockFilter struct {
	Calls int
	R     bool
	Err   error
	P     types.Provider
}

// Filter returns the specified mock result and/or error and increments calls.
func (m *MockFilter) Filter(p types.Provider) (bool, error) {
	m.P = p
	m.Calls++
	return m.R, m.Err
}
