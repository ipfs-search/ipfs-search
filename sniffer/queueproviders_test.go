package sniffer

import (
	"context"
	"errors"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestQueueContextCancel tests whether we're returning an error on context cancellation
func TestQueueContextCancel(t *testing.T) {
	assert := assert.New(t)

	q := &mockQueue{}

	c := make(chan types.Provider)

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context immediately
	cancel()

	pq := &providerQueuer{}
	err := pq.queue(ctx, c, q)

	assert.Equal(err, context.Canceled)
}

// TestQueuePublish tests whether a queued provider gets published.
func TestQueuePublish(t *testing.T) {
	assert := assert.New(t)

	pubs := make(chan interface{}, 1)
	priorities := make(chan uint8, 1)

	q := &mockQueue{
		pubs:       pubs,
		priorities: priorities,
	}

	p := mockProvider()
	c := make(chan types.Provider, 1)
	c <- p

	ctx, cancel := context.WithCancel(context.Background())

	pq := &providerQueuer{}
	go pq.queue(ctx, c, q)

	// See whether stuff actually got queued
	pub := <-pubs
	priority := <-priorities

	assert.Equal(pub.(*crawler.Args).Hash, p.Resource.ID)
	assert.Equal(priority, uint8(9))

	// Cleanup
	cancel()
}

// TestQueueError tests whether errors in publish are propagated
func TestQueueError(t *testing.T) {
	assert := assert.New(t)

	mockErr := errors.New("mock")

	pubs := make(chan interface{}, 1)
	priorities := make(chan uint8, 1)

	q := &mockQueue{
		pubs:       pubs,
		priorities: priorities,
		err:        mockErr,
	}

	p := mockProvider()
	c := make(chan types.Provider, 1)
	c <- p

	pq := &providerQueuer{}

	ctx := context.Background()
	err := pq.queue(ctx, c, q)

	assert.True(errors.Is(err, mockErr))
}
