package cache

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/ipfs-search/ipfs-search/components/index"
	"github.com/ipfs-search/ipfs-search/instr"
)

type testStruct struct {
	ValOne   string
	ValTwo   int
	ValThree float64
	ValFour  struct {
		v string
	}
	valNope int
}

type cacheStruct struct {
	ValOne   *string
	ValThree float64
	ValFour  struct {
		v string
	}
}

var cachingFields []string = []string{"ValOne", "ValThree", "ValFour", "neverFound"}
var props = testStruct{"test", 5, 5.5, struct{ v string }{"h"}, 3}
var testStr = "test"
var emptyStr = ""
var cachedProps = cacheStruct{&testStr, 5.5, struct{ v string }{"h"}}
var emptyCachedProps = cacheStruct{&emptyStr, 0.0, struct{ v string }{""}}

var testID = "testID"
var testErr = errors.New("errr")

type CacheTestSuite struct {
	suite.Suite
	ctx   context.Context
	instr *instr.Instrumentation

	i            *Index
	cachingIndex *index.Mock
	backingIndex *index.Mock
}

func (s *CacheTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.instr = instr.New()

	s.cachingIndex = &index.Mock{}
	s.backingIndex = &index.Mock{}

	s.i = &Index{
		backingIndex:    s.backingIndex,
		cachingIndex:    s.cachingIndex,
		cachingType:     reflect.TypeOf(cacheStruct{}),
		Instrumentation: s.instr,
	}
}

func (s *CacheTestSuite) TestNew() {
	// Allow value.
	i := New(s.backingIndex, s.cachingIndex, cacheStruct{}, s.instr)
	s.NotNil(i)

	// Allow pointer as well.
	i = New(s.backingIndex, s.cachingIndex, &cacheStruct{}, s.instr)
	s.NotNil(i)
}

func (s *CacheTestSuite) TestString() {
	exp := fmt.Sprintf("'%s' through '%s'", s.backingIndex, s.cachingIndex)
	s.Equal(exp, s.i.String())
}

func (s *CacheTestSuite) TestMakeCachingProperties() {
	res := s.i.makeCachingProperties(&props)

	s.NotNil(res)

	s.IsType(res, &cacheStruct{})

	cacheres := res.(*cacheStruct)
	s.Equal(props.ValOne, *cacheres.ValOne)
	s.Equal(props.ValThree, cacheres.ValThree)
	s.Equal(props.ValFour, cacheres.ValFour)
}

func (s *CacheTestSuite) TestIndexSuccess() {
	s.backingIndex.On("Index", mock.Anything, testID, &props).Return(nil).Once()
	s.cachingIndex.On("Index", mock.Anything, testID, &cachedProps).Return(nil).Once()

	err := s.i.Index(s.ctx, testID, &props)
	s.NoError(err)
}

func (s *CacheTestSuite) TestIndexBackingFail() {
	s.backingIndex.On("Index", mock.Anything, testID, &props).Return(testErr).Once()

	err := s.i.Index(s.ctx, testID, &props)
	s.Error(err)
	s.ErrorIs(err, testErr)

	// Backing fail: caching never called
	s.cachingIndex.AssertNumberOfCalls(s.T(), "Index", 0)

}

func (s *CacheTestSuite) TestIndexCachingFail() {
	s.backingIndex.On("Index", mock.Anything, testID, &props).Return(nil).Once()
	s.cachingIndex.On("Index", mock.Anything, testID, &cachedProps).Return(testErr).Once()

	err := s.i.Index(s.ctx, testID, &props)
	s.Error(err)
	s.ErrorIs(err, testErr)
	s.ErrorAs(err, &ErrCache{})
}

func (s *CacheTestSuite) TestUpdateSuccess() {
	s.backingIndex.On("Update", mock.Anything, testID, &props).Return(nil).Once()
	s.cachingIndex.On("Update", mock.Anything, testID, &cachedProps).Return(nil).Once()

	err := s.i.Update(s.ctx, testID, &props)
	s.NoError(err)
}

func (s *CacheTestSuite) TestUpdateBackingFail() {
	s.backingIndex.On("Update", mock.Anything, testID, &props).Return(testErr).Once()
	s.cachingIndex.On("Update", mock.Anything, testID, &cachedProps).Return(nil).Once()

	err := s.i.Update(s.ctx, testID, &props)
	s.Error(err)
	s.ErrorIs(err, testErr)

}

func (s *CacheTestSuite) TestUpdateCachingFail() {
	s.cachingIndex.On("Update", mock.Anything, testID, &cachedProps).Return(testErr).Once()

	err := s.i.Update(s.ctx, testID, &props)
	s.Error(err)
	s.ErrorIs(err, testErr)
	s.ErrorAs(err, &ErrCache{})

	// Cache yodate fail: backing never called.
	s.backingIndex.AssertNumberOfCalls(s.T(), "Index", 0)
}

func (s *CacheTestSuite) TestDeleteSuccess() {
	s.backingIndex.On("Delete", mock.Anything, testID).Return(nil).Once()
	s.cachingIndex.On("Delete", mock.Anything, testID).Return(nil).Once()

	err := s.i.Delete(s.ctx, testID)
	s.NoError(err)
}

func (s *CacheTestSuite) TestDeleteBackingFail() {
	s.backingIndex.On("Delete", mock.Anything, testID).Return(testErr).Once()
	s.cachingIndex.On("Delete", mock.Anything, testID).Return(nil).Once()

	err := s.i.Delete(s.ctx, testID)
	s.Error(err)
	s.ErrorIs(err, testErr)
}

func (s *CacheTestSuite) TestDeleteCachingFail() {
	s.cachingIndex.On("Delete", mock.Anything, testID).Return(testErr).Once()

	err := s.i.Delete(s.ctx, testID)
	s.Error(err)
	s.ErrorIs(err, testErr)
	s.ErrorAs(err, &ErrCache{})

	// Check: cache fails -> backing never called
	s.backingIndex.AssertNumberOfCalls(s.T(), "Delete", 0)
}

func (s *CacheTestSuite) TestGetCacheSuccess() {
	var data testStruct

	s.cachingIndex.On("Get", mock.Anything, testID, &data, mock.Anything).Return(true, nil).Once()

	found, err := s.i.Get(s.ctx, testID, &data)
	s.True(found)
	s.NoError(err)

	s.backingIndex.AssertNumberOfCalls(s.T(), "Get", 0)
}

func (s *CacheTestSuite) TestGetBackingSuccess() {
	var data testStruct

	s.cachingIndex.On("Get", mock.Anything, testID, &data, mock.Anything).Return(false, nil).Once()
	s.backingIndex.On("Get", mock.Anything, testID, &data, mock.Anything).Return(true, nil).Once()

	// If an item is not in the cache but is found in the backing, it will be added to the cache.
	s.cachingIndex.On("Index", mock.Anything, testID, &emptyCachedProps).Return(nil).Once()

	found, err := s.i.Get(s.ctx, testID, &data)
	s.True(found)
	s.NoError(err)
}

func (s *CacheTestSuite) TestGetCacheIndexFail() {
	var data testStruct

	s.cachingIndex.On("Get", mock.Anything, testID, &data, mock.Anything).Return(false, nil).Once()
	s.backingIndex.On("Get", mock.Anything, testID, &data, mock.Anything).Return(true, nil).Once()

	// When there is an error adding things to the cache, we should still get results
	s.cachingIndex.On("Index", mock.Anything, testID, &emptyCachedProps).Return(testErr).Once()

	found, err := s.i.Get(s.ctx, testID, &data)
	s.True(found)
	s.Error(err)
	s.ErrorIs(err, testErr)
	s.ErrorAs(err, &ErrCache{})
}

func (s *CacheTestSuite) TestGetCacheFail() {
	var data testStruct

	// If cache fails, we should continue to backing, so both get called.
	s.cachingIndex.On("Get", mock.Anything, testID, &data, mock.Anything).Return(false, testErr).Once()
	s.backingIndex.On("Get", mock.Anything, testID, &data, mock.Anything).Return(true, nil).Once()

	// We will still *try* to index it.
	s.cachingIndex.On("Index", mock.Anything, testID, &emptyCachedProps).Return(nil).Once()

	found, err := s.i.Get(s.ctx, testID, &data)
	s.True(found)
	s.Error(err)
	s.ErrorIs(err, testErr)
	s.ErrorAs(err, &ErrCache{})
}

func (s *CacheTestSuite) TestGetBackingFail() {
	var data testStruct

	s.cachingIndex.On("Get", mock.Anything, testID, &data, mock.Anything).Return(false, nil).Once()
	s.backingIndex.On("Get", mock.Anything, testID, &data, mock.Anything).Return(false, testErr).Once()

	found, err := s.i.Get(s.ctx, testID, &data)
	s.False(found)
	s.Error(err)
	s.ErrorIs(err, testErr)
}

func (s *CacheTestSuite) AfterTest() {
	s.cachingIndex.AssertExpectations(s.T())
	s.backingIndex.AssertExpectations(s.T())
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
