package index

import (
	"context"
	"errors"
	"log"

	"golang.org/x/sync/errgroup"
)

const debug bool = false

func contextDone(ctx context.Context, err error) bool {
	ctxErr := ctx.Err()
	if ctxErr != nil {
		return errors.Is(err, ctxErr)
	}

	return false
}

// MultiGet returns `fields` for the first document with `id` from given `indexes`.
// When the document is not found (nil, nil) is returned.
func MultiGet(ctx context.Context, indexes []Index, id string, dst interface{}, fields ...string) (Index, error) {
	foundIdx := make(chan Index, 1)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // cancel when we are finished

	g, groupCtx := errgroup.WithContext(ctx)
	for _, i := range indexes {
		i := i // https://go.dev/doc/faq#closures_and_goroutines

		g.Go(func() error {
			if debug {
				log.Printf("MultiGet %s index %s", id, i)
			}

			found, err := i.Get(groupCtx, id, dst, fields...)

			if err != nil && !contextDone(ctx, err) {
				// Ignore context done errors if MultiGet context is canceled.
				return err
			}

			if found {
				select {
				case <-groupCtx.Done():
					return nil
				case foundIdx <- i:
					cancel() // Found, we're done.
				}

			}

			return nil
		})
	}

	err := g.Wait()

	select {
	case result := <-foundIdx:
		return result, err
	default:
		return nil, err
	}
}
