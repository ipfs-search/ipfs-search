package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/ipfs-search/ipfs-search/instr"
)

type RedisTestSuite struct {
	suite.Suite
	ctx   context.Context
	instr *instr.Instrumentation

	i *Index
}

type rueidisMock struct {
	mock.Mock
}

func (s *RedisTestSuite) SetupTest() {

}

func (s *RedisTestSuite) TestGetKey() {

}

func (s *RedisTestSuite) TestSet() {

}

func (s *RedisTestSuite) TestString() {

}

func (s *RedisTestSuite) TestIndex() {

}

func (s *RedisTestSuite) TestUpdate() {

}

func (s *RedisTestSuite) TestDelete() {

}

func (s *RedisTestSuite) TestGet() {

}

func (s *RedisTestSuite) AfterTest() {
	// s.cachingIndex.AssertExpectations(s.T())
	// s.backingIndex.AssertExpectations(s.T())
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
