package worker

import (
	"context"
)

// WorkFunction is a function performing the work of a worker
type WorkFunction func(context.Context) error

// Function worker performs a single function as the process call
type Function struct {
	WorkFunc WorkFunction
}

// Work calls WorkFunc with context, returning any potential error
func (f *Function) Work(ctx context.Context) error {
	return f.WorkFunc(ctx)
}
