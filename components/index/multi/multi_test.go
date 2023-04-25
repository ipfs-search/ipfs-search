package multi

import (
	"context"
	"testing"

	"github.com/ipfs-search/ipfs-search/components/index"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type selectorMock struct {
	mock.Mock
}

func (m *selectorMock) ListIndexes() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *selectorMock) GetIndex(id string, properties Properties) string {
	args := m.Called(id, properties)
	return args.Get(0).(string)
}

// Compile-time assurance that implementation satisfies interface.
var _ Selector = &selectorMock{}

type MultiTestSuite struct {
	suite.Suite

	ctx      context.Context
	factory  index.FactoryMock
	selector selectorMock
	m        *Multi
	idx1     *index.Mock
	idx2     *index.Mock
}

func (s *MultiTestSuite) SetupTest() {
	s.idx1 = &index.Mock{}
	s.idx2 = &index.Mock{}
	s.ctx = context.Background()
	s.m = &Multi{
		selector: &s.selector,
		indexList: []index.Index{
			s.idx1, s.idx2,
		},
		indexMap: map[string]index.Index{
			"index1": s.idx1,
			"index2": s.idx2,
		},
	}
}

func (s *MultiTestSuite) TestNew() {
	s.selector.On("ListIndexes").Return([]string{}).Once()
	m := New(&s.factory, &s.selector)
	s.NotNil(m)
	s.selector.AssertExpectations(s.T())
}

func (s *MultiTestSuite) TestSetIndexes() {
	s.selector.On("ListIndexes").Return([]string{"index4", "index5"}).Once()
	s.factory.On("NewIndex", "index4").Return(s.idx2).Once()
	s.factory.On("NewIndex", "index5").Return(s.idx1).Once()
	s.m.setIndexes(&s.factory)

	s.Equal(s.idx2, s.m.indexMap["index4"])
	s.Equal(s.idx1, s.m.indexMap["index5"])
	s.Equal(s.idx2, s.m.indexList[0])
	s.Equal(s.idx1, s.m.indexList[1])

	s.selector.AssertExpectations(s.T())
	s.factory.AssertExpectations(s.T())
}

func (s *MultiTestSuite) TestGetIndexFound() {
	var dst interface{}

	s.idx1.On("Get", mock.Anything, "id", dst, []string{"metadata"}).Return(false, nil)
	s.idx2.On("Get", mock.Anything, "id", dst, []string{"metadata"}).Return(true, nil)

	found, err := s.m.Get(s.ctx, "id", dst, "metadata")
	s.True(found)
	s.NoError(err)

	s.idx1.AssertExpectations(s.T())
	s.idx2.AssertExpectations(s.T())
}

func (s *MultiTestSuite) TestGetIndexNotFound() {
	var dst interface{}

	s.idx1.On("Get", mock.Anything, "id", dst, []string{"metadata"}).Return(false, nil)
	s.idx2.On("Get", mock.Anything, "id", dst, []string{"metadata"}).Return(false, nil)

	found, err := s.m.Get(s.ctx, "id", dst, "metadata")
	s.False(found)
	s.NoError(err)

	s.idx1.AssertExpectations(s.T())
	s.idx2.AssertExpectations(s.T())
}

func (s *MultiTestSuite) TestIndex() {
	props := map[string]string{
		"hello": "world",
	}

	s.selector.On("GetIndex", "myid", props).Return("index2").Once()
	s.idx2.On("Index", mock.Anything, "myid", props).Return(nil).Once()

	err := s.m.Index(s.ctx, "myid", props)
	s.NoError(err)

	s.selector.AssertExpectations(s.T())
	s.idx2.AssertExpectations(s.T())
}

func (s *MultiTestSuite) TestUpdate() {
	props := map[string]string{
		"hello": "world",
	}

	s.selector.On("GetIndex", "myid", props).Return("index1").Once()
	s.idx1.On("Update", mock.Anything, "myid", props).Return(nil).Once()

	err := s.m.Update(s.ctx, "myid", props)
	s.NoError(err)

	s.selector.AssertExpectations(s.T())
	s.idx1.AssertExpectations(s.T())
}

func TestMultiTestSuite(t *testing.T) {
	suite.Run(t, new(MultiTestSuite))
}
