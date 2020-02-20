package sniffer

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Actual log message for testing
var msg = map[string]interface{}{
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

func TestResourceProvider(t *testing.T) {
	assert := assert.New(t)

	provider, err := Message(msg).ResourceProvider()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	assert.Equal(provider.Date, time.Date(2020, time.January, 21, 17, 28, 02, 501941007, time.UTC))
	assert.Equal(provider.Resource.Protocol, "ipfs")
	assert.Equal(provider.Resource.Id, "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp")
	assert.Equal(provider.Provider, "QmeTtFXm42Jb2todcKR538j6qHYxXt6suUzpF3rtT9FPSd")
}
