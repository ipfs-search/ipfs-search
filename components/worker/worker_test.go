package worker

import (
	"context"
	"fmt"
	"testing"

	"encoding/json"

	samqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
)

type crawlerMock struct {
	mock.Mock
}

func (m *crawlerMock) Crawl(ctx context.Context, r *t.AnnotatedResource) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

type loadLimiterMock struct {
	mock.Mock
}

func (m *loadLimiterMock) LoadLimit() error {
	args := m.Called()
	return args.Error(0)
}

type WorkerTestSuite struct {
	suite.Suite

	ll *loadLimiterMock
	c  *crawlerMock
	w  *Worker
}

func (s *WorkerTestSuite) SetupTest() {
	cfg := &Config{}
	s.c = &crawlerMock{}
	s.ll = &loadLimiterMock{}
	s.w = &Worker{
		cfg:             cfg,
		crawler:         s.c,
		ll:              s.ll,
		Instrumentation: instr.New(),
	}
}

// Test getResource(d samqp.Delivery) (*t.AnnotatedResource, error)
func (s *WorkerTestSuite) TestGetResourceValid() {
	// Test delivery containing JSON body of t.AnnotatedResource.

	// Set all fields.
	expected := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmHash",
		},
		Source: t.DirectorySource,
		Reference: t.Reference{
			Parent: &t.Resource{
				Protocol: t.IPFSProtocol,
				ID:       "QmOtherHash",
			},
			Name: "hello",
		},
	}

	// Convert expected to JSON.
	encoded, err := json.Marshal(expected)
	s.NoError(err)

	// Create delivery with JSON body.
	d := &samqp.Delivery{
		Body: encoded,
	}

	// Call getResource() with delivery.
	actual, err := getResource(d)

	// Validate response based on expected.
	s.NoError(err)
	s.Equal(expected, actual)
}

// Test getResource() with invalid input.
func (s *WorkerTestSuite) TestGetResourceInvalid() {
	invalidJSON := `invalidJSON`

	// Create delivery with invalid JSON body.
	d := &samqp.Delivery{
		Body: []byte(invalidJSON),
	}

	// Call getResource() with delivery.
	actual, err := getResource(d)

	// Validate response.
	s.Error(err)
	s.Nil(actual)
}

type mockAcknowledger struct {
	mock.Mock
}

func (m *mockAcknowledger) Ack(tag uint64, multiple bool) error {
	args := m.Called(tag, multiple)
	return args.Error(0)
}

func (m *mockAcknowledger) Nack(tag uint64, multiple bool, requeue bool) error {
	args := m.Called(tag, multiple, requeue)
	return args.Error(0)
}

func (m *mockAcknowledger) Reject(tag uint64, requeue bool) error {
	args := m.Called(tag, requeue)
	return args.Error(0)
}

// Test ackOrReject(err error, d *samqp.Delivery) error
func (s *WorkerTestSuite) TestAckOrRejectNoCrawlError() {
	// Test reject.

	acknowledger := &mockAcknowledger{}
	acknowledger.On("Ack", mock.Anything, false).Return(nil)

	// Create delivery.
	d := &samqp.Delivery{
		Acknowledger: acknowledger,
	}

	err := ackOrReject(nil, d)
	s.NoError(err)
	acknowledger.AssertExpectations(s.T())
}

func (s *WorkerTestSuite) TestAckOrRejectCrawlError() {
	// Test reject.

	acknowledger := &mockAcknowledger{}
	acknowledger.On("Reject", mock.Anything, false).Return(nil)

	// Create delivery.
	d := &samqp.Delivery{
		Acknowledger: acknowledger,
	}

	err := ackOrReject(fmt.Errorf("error"), d)
	s.NoError(err)
	acknowledger.AssertExpectations(s.T())
}

func (s *WorkerTestSuite) TestAckOrRejectRejectError() {
	// Error during reject propagates.

	acknowledger := &mockAcknowledger{}
	acknowledger.On("Reject", mock.Anything, false).Return(fmt.Errorf("error2"))

	// Create delivery.
	d := &samqp.Delivery{
		Acknowledger: acknowledger,
	}

	err := ackOrReject(fmt.Errorf("error"), d)
	s.Error(err)
	acknowledger.AssertExpectations(s.T())
}

func (s *WorkerTestSuite) TestAckOrRejectAckError() {
	// Error during Ack propagates.

	acknowledger := &mockAcknowledger{}
	acknowledger.On("Ack", mock.Anything, false).Return(fmt.Errorf("error2"))

	// Create delivery.
	d := &samqp.Delivery{
		Acknowledger: acknowledger,
	}

	err := ackOrReject(nil, d)
	s.Error(err)
	acknowledger.AssertExpectations(s.T())
}

func (s *WorkerTestSuite) TestCrawlDeliveryCrawlSuccess() {
	ctx := context.Background()

	acknowledger := &mockAcknowledger{}
	acknowledger.On("Ack", mock.Anything, false).Return(nil)

	s.ll.On("LoadLimit").Return(nil)

	// Create resource.
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmHash",
		},
		Source: t.DirectorySource,
		Reference: t.Reference{
			Parent: &t.Resource{
				Protocol: t.IPFSProtocol,
				ID:       "QmOtherHash",
			},
			Name: "hello",
		},
	}

	encoded, err := json.Marshal(r)
	s.NoError(err)

	// Create delivery.
	d := &samqp.Delivery{
		Acknowledger: acknowledger,
		Body:         encoded,
	}

	s.c.On("Crawl", mock.Anything, mock.Anything).Return(nil)

	err = s.w.crawlDelivery(ctx, d)
	s.NoError(err)
	s.c.AssertExpectations(s.T())
}

func (s *WorkerTestSuite) TestCrawlDeliveryAckOrRejectPropagate() {
	ctx := context.Background()

	acknowledger := &mockAcknowledger{}
	acknowledger.On("Ack", mock.Anything, false).Return(fmt.Errorf("error2"))

	s.ll.On("LoadLimit").Return(nil)

	// Create delivery.
	d := &samqp.Delivery{
		Acknowledger: acknowledger,
	}

	err := s.w.crawlDelivery(ctx, d)
	s.Error(err)
}

func TestWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(WorkerTestSuite))
}
