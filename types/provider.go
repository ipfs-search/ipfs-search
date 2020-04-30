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
