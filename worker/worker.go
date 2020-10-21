package worker

import (
	"context"
)

// Worker represents a single (blocking) worker
type Worker interface {
	Work(context.Context) error
}
