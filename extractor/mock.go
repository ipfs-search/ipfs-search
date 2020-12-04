package extractor

import (
	"context"
	"github.com/stretchr/testify/mock"

	t "github.com/ipfs-search/ipfs-search/types"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) Extract(ctx context.Context, r *t.AnnotatedResource, result interface{}) error {
	args := m.Called(ctx, r, result)
	return args.Error(0)
}

// Compile-time assurance that implementation satisfies interface.
var _ Extractor = &Mock{}
