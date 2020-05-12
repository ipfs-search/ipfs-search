package sniffer

import (
	"context"
	"errors"
	"github.com/ipfs-search/ipfs-search/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const longTime = time.Duration(time.Minute)
const shortTime = time.Duration(0)

// TestContextCancel tests whether we're returning an error on context cancellation
func TestContextCancel(t *testing.T) {
	l := &mockLogger{
		wait: time.Duration(longTime),
	}

	e := &mockExtractor{}

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan types.Provider)

	assert := assert.New(t)

	// Cancel context immediately
	cancel()

	err := yieldProviders(ctx, l, e, c, longTime)

	assert.Equal(err, context.Canceled)
}

// TestTimeout tests whether we're returning an error on timeout
func TestTimeout(t *testing.T) {
	l := &mockLogger{
		wait: time.Duration(longTime),
	}

	e := &mockExtractor{}

	ctx := context.Background()
	c := make(chan types.Provider)

	assert := assert.New(t)

	err := yieldProviders(ctx, l, e, c, shortTime)

	assert.Equal(ErrorLoggerTimeout, err)

}

// TestYieldProvider tests whether a provider is yielded for a provider message
func TestYieldProvider(t *testing.T) {
	mockProvider := types.Provider{
		Resource: &types.Resource{
			Protocol: "mockProtocol",
			Id:       "mockKey",
		},
	}

	// Send a nil message to assure a call to extractor
	l := newMockLogger()
	l.msgs <- nil

	e := &mockExtractor{
		provider: &mockProvider,
	}

	ctx := context.Background()
	c := make(chan types.Provider)

	assert := assert.New(t)

	go yieldProviders(ctx, l, e, c, longTime)

	provider := <-c
	assert.Equal(mockProvider, provider)
}

// TestLoggerError tests for error propagation from the logger
func TestLoggerError(t *testing.T) {
	errMock := errors.New("mock")

	// Create mock logger with associated messages
	errc := make(chan error, 1)
	errc <- errMock
	l := mockLogger{
		errc: errc,
	}

	e := &mockExtractor{}

	ctx := context.Background()
	c := make(chan types.Provider)

	assert := assert.New(t)

	err := yieldProviders(ctx, l, e, c, longTime)

	assert.Equal(err, errMock)
}

// TestProviderErr tests for error propagation from ResourceProvider
func TestProviderErr(t *testing.T) {
	errMock := errors.New("mock")

	l := newMockLogger()
	l.msgs <- nil

	e := &mockExtractor{
		err: errMock,
	}

	ctx := context.Background()
	c := make(chan types.Provider)

	assert := assert.New(t)

	err := yieldProviders(ctx, l, e, c, longTime)

	assert.Equal(err, errMock)
}
