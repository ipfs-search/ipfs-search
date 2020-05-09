package sniffer

import (
	"context"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/sniffer/filters"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNew does a burn test for New()
func TestNew(t *testing.T) {
	assert := assert.New(t)

	q := mockQueue{}

	cfg := DefaultConfig()

	s, e := New(cfg, q)

	assert.NotEmpty(s)
	assert.Empty(e)
}

// TestSniffCancel tests whether running Sniff() with a cancelled context returns with a context error.
func TestSniffCancel(t *testing.T) {
	assert := assert.New(t)

	l := mockLogger{}
	x := mockExtractor{}
	f := &filters.MockFilter{}

	cfg := DefaultConfig()

	s := &Sniffer{
		cfg:       cfg,
		filter:    f,
		extractor: x,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := s.Sniff(ctx, l)
	assert.Contains(err.Error(), "context canceled")
}

// TestLogToPublish tests the full chain from a log to a publish
func TestLogToPublish(t *testing.T) {
	assert := assert.New(t)

	// Create queue and channels to retreive published messages and priorities
	pubs := make(chan interface{})
	priorities := make(chan uint8)
	q := &mockQueue{
		pubs:       pubs,
		priorities: priorities,
	}

	// Create sniffer
	cfg := DefaultConfig()
	s, e := New(cfg, q)
	assert.NotEmpty(s)
	assert.Empty(e)

	// Create buffered message channel and send mock message
	msgs := make(chan map[string]interface{}, 1)
	mockMsg := map[string]interface{}{
		"Duration":     33190,
		"Logs":         []string{},
		"Operation":    "handleAddProvider",
		"ParentSpanID": 0,
		"SpanID":       6.999711555735423e+18,
		"Start":        "2020-01-21T17:28:02.501941007Z",
		"Tags": map[string]interface{}{
			"key":    "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
			"peer":   "QmeTtFXm42Jb2todcKR538j6qHYxXt6suUzpF3rtT9FPSd",
			"system": "dht",
		},
		"TraceID": 4.483443946463055e+18,
	}
	msgs <- mockMsg

	// Create mock logger with associated messages
	l := mockLogger{
		msgs: msgs,
	}

	// Create cancelable context for sniffer to work with
	ctx, cancel := context.WithCancel(context.Background())

	// Run sniffer in goroutine
	go s.Sniff(ctx, l)

	// Retreive publication
	pub := <-pubs
	priority := <-priorities

	assert.Equal(pub.(*crawler.Args).Hash, "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp")
	assert.Equal(priority, uint8(9))

	// Cleanup
	cancel()
}
