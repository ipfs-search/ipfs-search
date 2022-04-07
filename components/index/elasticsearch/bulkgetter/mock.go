package bulkgetter

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// Mock mocks the AsyncGetter interface.
type Mock struct {
	mock.Mock
}

// Get mocks a get of an AsyncGetter.
func (m *Mock) Get(ctx context.Context, req *GetRequest, dst interface{}) <-chan GetResponse {
	args := m.Called(ctx, req, dst)
	return args.Get(0).(<-chan GetResponse)
}

// Start mocks the start of an AsyncGetter.
func (m *Mock) Start(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Compile-time assurance that implementation satisfies interface.
var _ AsyncGetter = &Mock{}
