package queue

import (
	"context"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) Publish(pub interface{}, priority uint8) error {
	args := m.Called(pub, priority)
	return args.Error(0)
}

func (m *Mock) Consume() (<-chan amqp.Delivery, error) {
	args := m.Called()
	return args.Get(0).(<-chan amqp.Delivery), args.Error(1)
}

type MockFactory struct {
	mock.Mock
}

func (f *MockFactory) NewPublisher(ctx context.Context) (Publisher, error) {
	args := f.Called(ctx)
	return args.Get(0).(Publisher), args.Error(1)
}
