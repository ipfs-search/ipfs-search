package index

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MultiGetTestSuite struct {
	suite.Suite
	ctx  context.Context
	mock *Mock
}

func (s *MultiGetTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.mock = &Mock{}
	s.mock.Test(s.T())
}

// TestMultiGetNotFound tests "No document is found -> nil, 404 error"
func (s *MultiGetTestSuite) TestMultiGetNotFound() {
	dst := new(struct{})

	s.mock.On("Get", s.ctx, "objId", dst, []string{"testField"}).Return(false, nil)

	index, err := MultiGet(s.ctx, []Index{s.mock}, "objId", dst, "testField")

	s.Nil(index)
	s.NoError(err)
	s.mock.AssertExpectations(s.T())
}

// TestMultiGetFound tests "Document is found, with field not set"
func (s *MultiGetTestSuite) TestMultiGetFound() {
	dst := new(struct{})

	s.mock.On("Get", s.ctx, "objId", dst, []string{"testField"}).Return(true, nil)

	index, err := MultiGet(s.ctx, []Index{s.mock}, "objId", dst, "testField")

	s.NoError(err)
	s.Equal(index, s.mock)
	s.mock.AssertExpectations(s.T())
}

func TestMultiGetTestSuite(t *testing.T) {
	suite.Run(t, new(MultiGetTestSuite))
}
