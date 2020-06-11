package commands

import (
	"context"
	"log"
)

// block blocks until context is cancelled
func block(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

// log errors from errc
func errorLoop(errc <-chan error) {
	for {
		err := <-errc
		log.Printf("%T: %v", err, err)
	}
}
