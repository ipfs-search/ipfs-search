package index

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) Index(ctx context.Context, id string, properties interface{}) error {
	args := m.Called(ctx, id, properties)
	return args.Error(0)
}

func (m *Mock) Update(ctx context.Context, id string, properties interface{}) error {
	args := m.Called(ctx, id, properties)
	return args.Error(0)
}

func (m *Mock) Get(ctx context.Context, id string, dst interface{}, fields ...string) (bool, error) {
	args := m.Called(ctx, id, dst, fields)
	return args.Bool(0), args.Error(1)
}

// Compile-time assurance that implementation satisfies interface.
var _ Index = &Mock{}
