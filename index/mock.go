package index

import (
	"context"
	"github.com/stretchr/testify/mock"
	"reflect"
)

type mockResult struct {
	references []string
}

type mockIndex struct {
	mock.Mock

	ID     string
	Fields []string
	Result mockResult
	Found  bool
	Error  error
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
	// TODO; replace by generic mock

	m.ID = id
	m.Fields = fields

	// Set result
	v := reflect.ValueOf(dst).Elem()
	v.Set(reflect.ValueOf(m.Result))

	return m.Found, m.Error
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
