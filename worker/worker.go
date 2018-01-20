package worker

import (
	"golang.org/x/net/context"
)

// Worker represents a single (blocking) worker
type Worker interface {
	Work(context.Context) error
}

// Factory returns a worker
type Factory func() (Worker, error)
