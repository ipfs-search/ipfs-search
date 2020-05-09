package sniffer

import (
	"context"
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
