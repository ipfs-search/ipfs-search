package sniffer

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestNoHandleAddProvider tests for messages which are not relevant
func TestNoHandleAddProvider(t *testing.T) {
	// Actual log message for testing
	mockMsg := map[string]interface{}{
		"Operation": "other",
	}

	assert := assert.New(t)

	e := ProviderExtractor{}

	p, err := e.Extract(mockMsg)

	assert.Empty(err)
	assert.Empty(p)
}

func TestExtract(t *testing.T) {
	// Actual log message for testing
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

	assert := assert.New(t)

	e := ProviderExtractor{}

	p, err := e.Extract(mockMsg)
	assert.Empty(err)

	assert.Equal(p.Date, time.Date(2020, time.January, 21, 17, 28, 02, 501941007, time.UTC))
	assert.Equal(p.Resource.Protocol, "ipfs")
	assert.Equal(p.Resource.Id, "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp")
	assert.Equal(p.Provider, "QmeTtFXm42Jb2todcKR538j6qHYxXt6suUzpF3rtT9FPSd")
}

func TestExtractError(t *testing.T) {
	assert := assert.New(t)
	e := ProviderExtractor{}

	testError := func(msg map[string]interface{}, errContains string) {
		p, err := e.Extract(msg)

		assert.Empty(p)
		assert.Error(err)

		if errContains != "" {
			assert.Contains(err.Error(), errContains)
		}
	}

	// handleAddProvider should have all required fields, so this should error
	testError(map[string]interface{}{
		"Operation": "handleAddProvider",
	}, "")

	// invalid date
	testError(map[string]interface{}{
		"Duration":     33190,
		"Logs":         []string{},
		"Operation":    "handleAddProvider",
		"ParentSpanID": 0,
		"SpanID":       6.999711555735423e+18,
		"Start":        "invalid date",
		"Tags": map[string]interface{}{
			"key":    "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
			"peer":   "QmeTtFXm42Jb2todcKR538j6qHYxXt6suUzpF3rtT9FPSd",
			"system": "dht",
		},
		"TraceID": 4.483443946463055e+18,
	}, "Error converting 'Start' into time")

	// missing date
	testError(map[string]interface{}{
		"Duration":     33190,
		"Logs":         []string{},
		"Operation":    "handleAddProvider",
		"ParentSpanID": 0,
		"SpanID":       6.999711555735423e+18,
		"Tags": map[string]interface{}{
			"key":    "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
			"peer":   "QmeTtFXm42Jb2todcKR538j6qHYxXt6suUzpF3rtT9FPSd",
			"system": "dht",
		},
		"TraceID": 4.483443946463055e+18,
	}, "'Start' not found in message")

	// missing tags
	testError(map[string]interface{}{
		"Duration":     33190,
		"Logs":         []string{},
		"Operation":    "handleAddProvider",
		"ParentSpanID": 0,
		"SpanID":       6.999711555735423e+18,
		"Start":        "2020-01-21T17:28:02.501941007Z",
		"TraceID":      4.483443946463055e+18,
	}, "'Tags' not found in message:")

	// missing tag key
	testError(map[string]interface{}{
		"Duration":     33190,
		"Logs":         []string{},
		"Operation":    "handleAddProvider",
		"ParentSpanID": 0,
		"SpanID":       6.999711555735423e+18,
		"Start":        "2020-01-21T17:28:02.501941007Z",
		"Tags": map[string]interface{}{
			"peer":   "QmeTtFXm42Jb2todcKR538j6qHYxXt6suUzpF3rtT9FPSd",
			"system": "dht",
		},
		"TraceID": 4.483443946463055e+18,
	}, "Could not read 'key' in tags of message:")

	// missing tag peer
	testError(map[string]interface{}{
		"Duration":     33190,
		"Logs":         []string{},
		"Operation":    "handleAddProvider",
		"ParentSpanID": 0,
		"SpanID":       6.999711555735423e+18,
		"Start":        "2020-01-21T17:28:02.501941007Z",
		"Tags": map[string]interface{}{
			"key":    "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
			"system": "dht",
		},
		"TraceID": 4.483443946463055e+18,
	}, "Could not read 'peer' in tags of message:")
}
