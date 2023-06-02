package index

import (
	"github.com/stretchr/testify/mock"
)

// FactoryMock mocks a Factory, returning an Index Mock.
type FactoryMock struct {
	mock.Mock
}

// NewIndex returns an Index Mock.
func (m *FactoryMock) NewIndex(name string) Index {
	m.Called(name)
	return &Mock{}
}

// Compile-time assurance that implementation satisfies interface.
var _ Factory = &FactoryMock{}
