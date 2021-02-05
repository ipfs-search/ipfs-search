package protocol

import (
	"context"
	"github.com/stretchr/testify/mock"

	t "github.com/ipfs-search/ipfs-search/types"
)

// Mock mocks the Protocol interface.
type Mock struct {
	mock.Mock
}

// GatewayURL mocks the corresponding method on the Protocol interface.
func (m *Mock) GatewayURL(r *t.AnnotatedResource) string {
	args := m.Called(r)
	return args.String(0)
}

// Stat mocks the corresponding method on the Protocol interface.
func (m *Mock) Stat(ctx context.Context, r *t.AnnotatedResource) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

// Ls mocks the corresponding method on the Protocol interface.
func (m *Mock) Ls(ctx context.Context, r *t.AnnotatedResource, c chan<- *t.AnnotatedResource) error {
	args := m.Called(ctx, r, c)
	return args.Error(0)
}

// IsInvalidResourceErr mocks the corresponding method on the Protocol interface.
func (m *Mock) IsInvalidResourceErr(err error) bool {
	args := m.Called(err)
	return args.Bool(0)
}

// Compile-time assurance that implementation satisfies interface.
var _ Protocol = &Mock{}
