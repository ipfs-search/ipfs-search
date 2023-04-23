package aliasresolver

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// Mock mocks the Index interface.
type Mock struct {
	mock.Mock
}

func (m *Mock) GetIndex(ctx context.Context, aliasName string) (string, error) {
	args := m.Called(ctx, aliasName)
	return args.String(0), args.Error(1)
}

func (m *Mock) GetAlias(ctx context.Context, indexName string) (string, error) {
	args := m.Called(ctx, indexName)
	return args.String(0), args.Error(1)
}

// Compile-time assurance that implementation satisfies interface.
var _ AliasResolver = &Mock{}
