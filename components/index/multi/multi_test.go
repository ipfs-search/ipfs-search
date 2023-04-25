package multi

import (
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

	factory  index.FactoryMock
	selector selectorMock
	m        *Multi
}

func (s *MultiTestSuite) SetupTest() {
	s.m = &Multi{
		selector: &s.selector,
	}
}

func (s *MultiTestSuite) TestNew() {
	s.selector.On("ListIndexes").Return([]string{}).Once()
	m := New(&s.factory, &s.selector)
	s.NotNil(m)
	s.selector.AssertExpectations(s.T())
}

func (s *MultiTestSuite) TestSetIndexes() {
	s.selector.On("ListIndexes").Return([]string{"index1", "index2"}).Once()
	s.factory.On("NewIndex", "index1").Once()
	s.factory.On("NewIndex", "index2").Once()
	s.m.setIndexes(&s.factory)

	s.selector.AssertExpectations(s.T())
	s.factory.AssertExpectations(s.T())
}

func (s *MultiTestSuite) TestGetIndex() {

}

func (s *MultiTestSuite) TestIndex() {

}

func (s *MultiTestSuite) TestUpdate() {
}

func (s *MultiTestSuite) TestGet() {
}

func TestMultiTestSuite(t *testing.T) {
	suite.Run(t, new(MultiTestSuite))
}
