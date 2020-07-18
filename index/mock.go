package index

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type mockIndex struct {
	mock.Mock
}

func (m *mockIndex) Index(ctx context.Context, id string, properties map[string]interface{}) error {
	args := m.Called(ctx, id, properties)
	return args.Error(0)
}

func (m *mockIndex) Update(ctx context.Context, id string, properties map[string]interface{}) error {
	args := m.Called(ctx, id, properties)
	return args.Error(0)
}

func (m *mockIndex) Get(ctx context.Context, id string, dst interface{}, fields ...string) (bool, error) {
	args := m.Called(ctx, id, dst, fields)
	return args.Bool(0), args.Error(1)
}

func (m *mockIndex) Exists(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *mockIndex) Create(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockIndex) ConfigUpToDate(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *mockIndex) ConfigUpdate(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
