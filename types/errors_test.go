package types

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProviderErrorf(t *testing.T) {
	assert := assert.New(t)

	wrapped := errors.New("test")

	mockProvider := Provider{
		Resource: &Resource{
			Protocol: "mockProtocol",
			ID:       "mockKey",
		},
	}

	err := NewProviderErrorf(wrapped, mockProvider, "Format %s", "me")
	assert.Error(err)

	// Test wrapping
	assert.True(errors.Is(err, wrapped))

	// Test formatting
	assert.Equal(err.Error(), "Format me")
}
