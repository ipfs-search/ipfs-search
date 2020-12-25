package queuer

import (
	"context"
	"errors"
	"testing"

	"github.com/ipfs-search/ipfs-search/queue"
	t "github.com/ipfs-search/ipfs-search/types"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type QueuerTestSuite struct {
	suite.Suite
	q      *queue.Mock
	ctx    context.Context
	cancel func()
	p      t.Provider
	r      *t.AnnotatedResource
}

func (s *QueuerTestSuite) SetupTest() {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.q = &queue.Mock{}
	s.q.Test(s.T())
	s.p = t.MockProvider()
	s.r = &t.AnnotatedResource{
		Resource: s.p.Resource,
	}
}

func (s *QueuerTestSuite) TearDownTest() {
	s.cancel()
}

// TestQueueContextCancel tests whether we're returning an error on context cancellation
func (s *QueuerTestSuite) TestQueueContextCancel() {
	ch := make(chan t.Provider)

	// Cancel context immediately
	s.cancel()

	pq := New(s.q, ch)

	err := pq.Queue(s.ctx)

	s.Equal(err, context.Canceled)

	s.q.AssertNotCalled(s.T(), "Publish")
}

// TestQueuePublish tests whether a queued provider gets published.
func (s *QueuerTestSuite) TestQueuePublish() {
	s.q.On("Publish", mock.Anything, s.r, uint8(9)).Return(nil)

	ch := make(chan t.Provider)

	go func() {
		// Process provider
		ch <- s.p

		// Cancel when done, causing pq to return
		s.cancel()
	}()

	pq := New(s.q, ch)
	err := pq.Queue(s.ctx)

	s.Equal(err, context.Canceled)

	s.q.AssertExpectations(s.T())
}

// TestQueueError tests whether errors in publish are propagated
func (s *QueuerTestSuite) TestQueueError() {
	mockErr := errors.New("mock")

	s.q.On("Publish", mock.Anything, s.r, uint8(9)).Return(mockErr)

	ch := make(chan t.Provider)

	go func() {
		// Process provider
		ch <- s.p

		// Cancel when done, causing pq to return`
		s.cancel()
	}()

	pq := New(s.q, ch)
	err := pq.Queue(s.ctx)

	s.True(errors.Is(err, mockErr))

	s.q.AssertExpectations(s.T())
}

func TestQueuerTestSuite(t *testing.T) {
	suite.Run(t, new(QueuerTestSuite))
}
