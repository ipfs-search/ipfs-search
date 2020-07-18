package index

import (
	"context"
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

type EnsureTestSuite struct {
	suite.Suite
	ctx  context.Context
	mock *mockIndex
}

func (s *EnsureTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.mock = &mockIndex{}
	s.mock.Test(s.T())
}

func (s *EnsureTestSuite) TestEnsureExistsCreate() {
	s.mock.On("Exists", s.ctx).Return(false, nil).Once() // Doesn't exist yet
	s.mock.On("Create", s.ctx).Return(nil).Once()

	err := ensureExists(s.ctx, s.mock)

	s.NoError(err)
	s.mock.AssertExpectations(s.T())
}

func (s *EnsureTestSuite) TestEnsureExistsCreateExistsError() {
	existsErr := errors.New("existsErr")

	s.mock.On("Exists", s.ctx).Return(false, existsErr).Once()

	err := ensureExists(s.ctx, s.mock)

	s.True(errors.Is(err, existsErr))
	s.mock.AssertExpectations(s.T())
}

func (s *EnsureTestSuite) TestEnsureExistsCreateCreateError() {
	createErr := errors.New("createErr")

	s.mock.On("Exists", s.ctx).Return(false, nil).Once() // Doesn't exist yet
	s.mock.On("Create", s.ctx).Return(createErr).Once()

	err := ensureExists(s.ctx, s.mock)

	s.True(errors.Is(err, createErr))
	s.mock.AssertExpectations(s.T())
}

func (s *EnsureTestSuite) TestEnsureConfigUpToDateUpToDateErr() {
	upToDateErr := errors.New("upToDateErr")

	s.mock.On("ConfigUpToDate", s.ctx).Return(false, upToDateErr).Once() // Initially: not up to date

	err := ensureConfigUpToDate(s.ctx, s.mock)

	s.True(errors.Is(err, upToDateErr))
	s.mock.AssertExpectations(s.T())
}

func (s *EnsureTestSuite) TestEnsureConfigUpToDateUpdateErr() {
	updateErr := errors.New("updateErr")

	s.mock.On("ConfigUpToDate", s.ctx).Return(false, nil).Once() // Initially: not up to date
	s.mock.On("ConfigUpdate", s.ctx).Return(updateErr).Once()    // Simulate update, no error

	err := ensureConfigUpToDate(s.ctx, s.mock)

	s.True(errors.Is(err, updateErr))
	s.mock.AssertExpectations(s.T())
}

func (s *EnsureTestSuite) TestEnsureConfigUpToDateUpdateSuccess() {
	s.mock.On("ConfigUpToDate", s.ctx).Return(false, nil).Once() // Initially: not up to date
	s.mock.On("ConfigUpdate", s.ctx).Return(nil).Once()          // Simulate update, no error
	s.mock.On("ConfigUpToDate", s.ctx).Return(true, nil).Once()  // After: up to date

	err := ensureConfigUpToDate(s.ctx, s.mock)

	s.NoError(err)
	s.mock.AssertExpectations(s.T())
}

func (s *EnsureTestSuite) TestEnsureConfigUpToDateUpdateFail() {
	s.mock.On("ConfigUpToDate", s.ctx).Return(false, nil).Once() // Initially: not up to date
	s.mock.On("ConfigUpdate", s.ctx).Return(nil).Once()          // Simulate update, no error
	s.mock.On("ConfigUpToDate", s.ctx).Return(false, nil).Once() // After: not updated

	err := ensureConfigUpToDate(s.ctx, s.mock)

	s.Error(err)
	s.mock.AssertExpectations(s.T())
}

func (s *EnsureTestSuite) TestEnsureExistsAndUpdated() {
	s.mock.On("Exists", s.ctx).Return(true, nil).Once()         // Exists
	s.mock.On("ConfigUpToDate", s.ctx).Return(true, nil).Once() // Up to date

	err := ensureExistsAndUpdated(s.ctx, s.mock)

	s.NoError(err)
	s.mock.AssertExpectations(s.T())
}

func (s *EnsureTestSuite) TestEnsureExistsAndUpdatedMulti() {
	s.mock.On("Exists", s.ctx).Return(true, nil).Twice()         // Exists
	s.mock.On("ConfigUpToDate", s.ctx).Return(true, nil).Twice() // Up to date

	err := EnsureExistsAndUpdatedMulti(s.ctx, s.mock, s.mock)

	s.NoError(err)
	s.mock.AssertExpectations(s.T())
}

func TestEnsureTestSuite(t *testing.T) {
	suite.Run(t, new(EnsureTestSuite))
}
