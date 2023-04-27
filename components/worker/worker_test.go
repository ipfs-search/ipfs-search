package worker

import (
	"context"
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

func TestWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(WorkerTestSuite))
}
