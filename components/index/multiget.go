package index

import (
	// "log"
	"context"

	"golang.org/x/sync/errgroup"
)

// MultiGet returns `fields` for the first document with `id` from given `indexes`.
// When the document is not found (nil, nil) is returned.
func MultiGet(ctx context.Context, indexes []Index, id string, dst interface{}, fields ...string) (Index, error) {
	var foundIdx Index

	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // cancel when we are finished

	g, groupCtx := errgroup.WithContext(ctx)
	for _, i := range indexes {
		i := i // https://go.dev/doc/faq#closures_and_goroutines

		g.Go(func() error {
			// log.Printf("MultiGet %s index %s", id, i)

			found, err := i.Get(groupCtx, id, dst, fields...)

			if err != nil {
				return err
			}

			if found {
				cancel() // We're done
				foundIdx = i
			}

			return nil
		})
	}

	return foundIdx, g.Wait()
}
