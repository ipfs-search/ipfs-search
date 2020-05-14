package index

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestMultiGetNotFound tests "No document is found -> nil, 404 error"
func TestMultiGetNotFound(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	m := &mockIndex{
		Found: false,
	}

	var dst interface{}

	index, err := MultiGet(ctx, []Index{m}, "", &dst, "")

	assert.Nil(index)
	assert.Nil(err)
}

// TestMultiGetFound tests "Document is found, with field not set"
func TestMultiGetFound(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	// Container for query reference fetch results
	dst := new(mockResult)
	res := []string{"hoi", "doei"}

	m := &mockIndex{
		Found: true,
		Result: mockResult{
			references: res,
		},
	}

	index, err := MultiGet(ctx, []Index{m}, "", dst, "references")

	assert.Equal(index, m)
	assert.Equal(res, dst.references)
	assert.Nil(err)
}
