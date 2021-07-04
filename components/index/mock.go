package index

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// Mock mocks the Index interface.
type Mock struct {
	mock.Mock
}

// Index mocks the Index method on the Index interface.
func (m *Mock) Index(ctx context.Context, id string, properties interface{}) error {
	args := m.Called(ctx, id, properties)
	return args.Error(0)
}

// Update mocks the Update method on the Index interface.
func (m *Mock) Update(ctx context.Context, id string, properties interface{}) error {
	args := m.Called(ctx, id, properties)
	return args.Error(0)
}

// Get mocks the Get method on the Index interface.
func (m *Mock) Get(ctx context.Context, id string, dst interface{}, fields ...string) (bool, error) {
	args := m.Called(ctx, id, dst, fields)
	return args.Bool(0), args.Error(1)
}

// Delete mocks the Delete method on the Index interface.
func (m *Mock) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Compile-time assurance that implementation satisfies interface.
var _ Index = &Mock{}
