package crawler

import (
	"context"
)

// Queue allows publishing of sniffed items.
type Queue interface {
	Publish(context.Context, interface{}, uint8) error
}
