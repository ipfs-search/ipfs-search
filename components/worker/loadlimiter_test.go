package worker

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const (
	tMin  = time.Millisecond
	tMax  = 10 * time.Millisecond
	lLow  = 0.1
	lMax  = 0.8
	lHigh = 0.9
)

type loadMock struct {
	mock.Mock
}

func (m *loadMock) getLoad() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

func (m *loadMock) throttle(t time.Duration) {
	m.Called(t)
	return
}

type LoadLimiterTestSuite struct {
	suite.Suite

	l *loadLimiter
	m *loadMock
}

func (s *LoadLimiterTestSuite) SetupTest() {
	s.m = &loadMock{}

	s.l = &loadLimiter{
		name:            "test",
		maxLoad:         lMax,
		throttleCurrent: tMin,
		throttleMin:     tMin,
		throttleMax:     tMax,
		throttle:        s.m.throttle,
		getLoad:         s.m.getLoad,
	}
}

func (s *LoadLimiterTestSuite) TestIntegrationLowLoad() {
	ll := NewLoadLimiter("test", 5.0, 100*tMin, 100*tMax)
	before := time.Now()
	s.NoError(ll.LoadLimit())
	after := time.Now()

	// Allow 10ms deviation.
	s.WithinDuration(before, after, 10*tMin)
}

func (s *LoadLimiterTestSuite) TestIntegrationHighLoad() {
	// Assume some load on system, this should cause a 100ms delay.
	ll := NewLoadLimiter("test", 0.00, 100*tMin, 100*tMax)
	before := time.Now()
	s.NoError(ll.LoadLimit())
	after := time.Now()

	expected := before.Add(100 * tMin)

	// Allow 10ms deviation + 100*tMin jitter.
	s.WithinDuration(expected, after, 210*tMin)
}

func (s *LoadLimiterTestSuite) TestThrottleDouble() {
	// First time doubles.
	s.l.throttleDouble()
	s.Equal(2*tMin, s.l.throttleCurrent)

	// Ensure maximum is kept.
	for i := 0; i < 100; i++ {
		s.l.throttleDouble()
	}

	s.Equal(tMax, s.l.throttleCurrent)
}

func (s *LoadLimiterTestSuite) TestThrottleReset() {
	s.l.throttleReset()
	s.Equal(tMin, s.l.throttleCurrent)

	s.l.throttleDouble()
	s.l.throttleReset()
	s.Equal(tMin, s.l.throttleCurrent)

	for i := 0; i < 100; i++ {
		s.l.throttleDouble()
	}
	s.l.throttleReset()
	s.Equal(tMin, s.l.throttleCurrent)
}

func (s *LoadLimiterTestSuite) TestLoadLimitLowLoad() {
	s.m.On("getLoad").Once().Return(lLow, nil)

	err := s.l.LoadLimit()
	s.NoError(err)
}

func (s *LoadLimiterTestSuite) TestLoadLimitHighLoad() {
	s.m.On("getLoad").Return(lHigh, nil)
	s.m.On("throttle", tMin)

	err := s.l.LoadLimit()
	s.NoError(err)

	s.m.AssertNumberOfCalls(s.T(), "getLoad", 1)
	s.m.AssertNumberOfCalls(s.T(), "throttle", 1)
}

func (s *LoadLimiterTestSuite) TestLoadLimitReset() {
	s.m.On("getLoad").Once().Return(lHigh, nil)
	s.m.On("throttle", mock.Anything)

	err := s.l.LoadLimit()
	s.NoError(err)

	s.m.On("getLoad").Once().Return(lLow, nil)
	err = s.l.LoadLimit()
	s.NoError(err)

	s.Equal(tMin, s.l.throttleCurrent)
}

func (s *LoadLimiterTestSuite) TestLoadLimitError() {
	s.m.On("getLoad").Return(0.0, errors.New("test"))
	err := s.l.LoadLimit()
	s.Error(err, "test")
}

func (s *LoadLimiterTestSuite) AfterTest() {
	s.m.AssertExpectations(s.T())
}

func TestLoadLimiterTestSuite(t *testing.T) {
	suite.Run(t, new(LoadLimiterTestSuite))
}
