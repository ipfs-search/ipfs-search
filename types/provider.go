package types

import (
	"fmt"
	"time"
)

// Provider represents a Resource available from an identified provider at a particular moment.
type Provider struct {
	*Resource
	Date     time.Time
	Provider string
}

// String defaults to the URI
func (r *Provider) String() string {
	return fmt.Sprintf("%s at %s on %s", r.URI(), r.Provider, r.Date)
}

// MockProvider returns a provider to be used for mocking.
func MockProvider() Provider {
	resource := &Resource{
		Protocol: "ipfs",
		ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
	}

	return Provider{
		Resource: resource,
		Date:     time.Now(),
		Provider: "QmeTtFXm42Jb2todcKR538j6qHYxXt6suUzpF3rtT9FPSd",
	}
}
