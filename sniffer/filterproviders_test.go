package sniffer

import (
	"context"
	"errors"
	"fmt"
	"github.com/ipfs-search/ipfs-search/sniffer/filters"
	"github.com/ipfs-search/ipfs-search/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterProvidersInclude(t *testing.T) {
	assert := assert.New(t)

	// Pass anything
	mockFilter := &filters.MockFilter{
		R: true,
	}

	p := types.Provider{}

	in := make(chan types.Provider)
	out := make(chan types.Provider)
	ctx, cancel := context.WithCancel(context.Background())

	go filterProviders(ctx, in, out, mockFilter)

	// Write provider
	in <- p

	// Read result
	r := <-out
	assert.Equal(p, r)

	// Make sure queue is empty
	assert.Equal(0, len(out))

	// Teardown goroutine
	cancel()
}

func TestFilterProvidersExclude(t *testing.T) {
	assert := assert.New(t)

	// Fail anything
	mockFilter := &filters.MockFilter{
		R: false,
	}

	p := types.Provider{}

	in := make(chan types.Provider)
	out := make(chan types.Provider)
	ctx, cancel := context.WithCancel(context.Background())

	go filterProviders(ctx, in, out, mockFilter)

	// Write provider
	in <- p

	// Make sure queue is empty
	assert.Equal(0, len(out))

	// Teardown goroutine
	cancel()
}

func TestFilterProvidersError(t *testing.T) {
	assert := assert.New(t)

	mockError := errors.New("test")
	mockFilter := &filters.MockFilter{
		Err: mockError,
	}

	p := types.Provider{}

	// Note the buffer, so we can write before calling filterProviders
	in := make(chan types.Provider, 1)
	out := make(chan types.Provider)
	ctx := context.Background()

	// Write provider
	in <- p

	err := filterProviders(ctx, in, out, mockFilter)

	assert.Error(err)
	assert.Equal(mockError, err)

	// Make sure queue is empty
	assert.Equal(0, len(out))
}

func TestFilterProvidersContextCancel(t *testing.T) {
	assert := assert.New(t)

	mockError := errors.New("test")
	mockFilter := &filters.MockFilter{
		Err: mockError,
	}

	// Note the buffer, so we can write before calling filterProviders
	in := make(chan types.Provider, 1)
	out := make(chan types.Provider)
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel already
	cancel()

	err := filterProviders(ctx, in, out, mockFilter)

	assert.Error(err)

	// As no providers are fed to the channel, no filters should be applied
	// so no error should propagate from there
	assert.Equal(mockFilter.Calls, 0, fmt.Sprintf("mockFilter has been called with provider %v", mockFilter.P))
	assert.Contains(err.Error(), "context canceled")

	// Make sure queue is empty
	assert.Equal(0, len(out))
}
