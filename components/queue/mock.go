package queue

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/mock"
)

// Mock mocks the Queue interface.
type Mock struct {
	mock.Mock
}

// Publish mocks the corresponding method on the Queue interface.
func (m *Mock) Publish(ctx context.Context, pub interface{}, priority uint8) error {
	args := m.Called(ctx, pub, priority)
	return args.Error(0)
}

// Consume mocks the corresponding method on the Queue interface.
func (m *Mock) Consume(ctx context.Context) (<-chan amqp.Delivery, error) {
	args := m.Called(ctx)
	return args.Get(0).(<-chan amqp.Delivery), args.Error(1)
}

// MockFactory mocks the Factory interface.
type MockFactory struct {
	mock.Mock
}

// NewPublisher mocks the corresponding method on the Factory interface.
func (f *MockFactory) NewPublisher(ctx context.Context) (Publisher, error) {
	args := f.Called(ctx)
	return args.Get(0).(Publisher), args.Error(1)
}

// Compile-time assurance that implementation satisfies interface.
var _ Queue = &Mock{}
var _ PublisherFactory = &MockFactory{}
