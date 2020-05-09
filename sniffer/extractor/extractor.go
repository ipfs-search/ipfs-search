package extractor

import (
	"fmt"
	t "github.com/ipfs-search/ipfs-search/types"
	"time"
)

// Extractor attempts to extract a ResourceProvider from a Message, returning nil when
// none was found and an error in unexpected situations.
type Extractor struct{}

// New returns a pointer to a new Extractor
func New() (*Extractor, error) {
	return &Extractor{}, nil
}

// Extract performs the extraction.
func (e Extractor) Extract(m map[string]interface{}) (*t.Provider, error) {
	// Somehow, real life messages are divided into events and operations.
	// This is not properly documented anywhere.
	operationType := m["Operation"]
	if operationType == "handleAddProvider" {
		rawDate, ok := m["Start"]
		if !ok {
			return nil, fmt.Errorf("'Start' not found in message: %v", m)
		}

		date, err := time.Parse(time.RFC3339Nano, rawDate.(string))
		if err != nil {
			return nil, fmt.Errorf("Error converting 'Start' into time: %w", err)
		}

		rawTags, ok := m["Tags"]
		if !ok {
			return nil, fmt.Errorf("'Tags' not found in message: %#v", m)
		}

		tags, ok := rawTags.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Could not convert 'Tags' for message: %#v", m)
		}

		key, ok := tags["key"].(string)
		if !ok {
			return nil, fmt.Errorf("Could not read 'key' in tags of message: %#v", m)
		}

		peer, ok := tags["peer"].(string)
		if !ok {
			return nil, fmt.Errorf("Could not read 'peer' in tags of message: %#v", m)
		}

		return &t.Provider{
			Resource: &t.Resource{
				Protocol: "ipfs",
				Id:       key,
			},
			Date:     date,
			Provider: peer,
		}, nil
	}

	return nil, nil
}
