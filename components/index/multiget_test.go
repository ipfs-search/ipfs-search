package index

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MultiGetTestSuite struct {
	suite.Suite
	ctx context.Context

	mock1 *Mock
	mock2 *Mock

	indexes []Index
}

func (s *MultiGetTestSuite) SetupTest() {
	s.ctx = context.Background()

	s.mock1 = &Mock{}
	s.mock1.Test(s.T())
	s.mock2 = &Mock{}
	s.mock2.Test(s.T())

	s.indexes = []Index{s.mock1, s.mock2}
}

func (s *MultiGetTestSuite) AfterTest() {
	s.mock1.AssertExpectations(s.T())
	s.mock2.AssertExpectations(s.T())
}

// TestMultiGetNotFound tests "No document is found -> nil, 404 error"
func (s *MultiGetTestSuite) TestNotFound() {
	dst := new(struct{})

	s.mock1.On("Get", mock.Anything, "objId", dst, []string{"testField"}).Return(false, nil)
	s.mock2.On("Get", mock.Anything, "objId", dst, []string{"testField"}).Return(false, nil)

	index, err := MultiGet(s.ctx, s.indexes, "objId", dst, "testField")

	s.Nil(index)
	s.NoError(err)
}

// TestMultiGetFound tests "Document is found, with field not set"
func (s *MultiGetTestSuite) TestFound() {
	type dstStruct struct {
		Changed bool
	}

	dst := dstStruct{false}

	s.mock1.On("Get", mock.Anything, "objId", &dst, []string{"testField"}).Run(func(args mock.Arguments) {
		u := args.Get(2).(*dstStruct)
		u.Changed = true
	}).Return(true, nil)
	s.mock2.On("Get", mock.Anything, "objId", &dst, []string{"testField"}).Return(false, nil)

	index, err := MultiGet(s.ctx, s.indexes, "objId", &dst, "testField")

	s.NoError(err)
	s.Equal(index, s.mock1)
	s.True(dst.Changed)
}

// TestMultiFound tets for predictable behaviour in case the item is found in multiple indexes.
// This implies a problem and hence should return an error.
func (s *MultiGetTestSuite) TestMultiFound() {
	type dstStruct struct {
		Value int
	}

	dst := dstStruct{}

	s.mock1.On("Get", mock.Anything, "objId", mock.Anything, []string{"testField"}).Run(func(args mock.Arguments) {
		u := args.Get(2).(*dstStruct)
		u.Value = 1
	}).Return(true, nil)
	s.mock2.On("Get", mock.Anything, "objId", mock.Anything, []string{"testField"}).Run(func(args mock.Arguments) {
		u := args.Get(2).(*dstStruct)
		u.Value = 2
	}).Return(true, nil)

	index, err := MultiGet(s.ctx, s.indexes, "objId", &dst, "testField")

	s.NoError(err)

	// This returns *either* one or the other.
	s.True(dst.Value == 1 || dst.Value == 2)
	s.True(index == s.mock1 || index == s.mock2)
}

func TestMultiGetTestSuite(t *testing.T) {
	suite.Run(t, new(MultiGetTestSuite))
}
