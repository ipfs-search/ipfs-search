package sniffer

import (
	"fmt"
	"time"
)

// Provider for a Resource at a particular Date.
type Provider struct {
	*Resource
	Date     time.Time
	Provider string
}

// String defaults to the URI
func (r *Provider) String() string {
	return fmt.Sprintf("%s at %s on %s", r.URI(), r.Provider, r.Date)
}
