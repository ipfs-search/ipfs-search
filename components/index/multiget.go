package index

import (
	"context"
)

// MultiGet returns `fields` for the first document with `id` from given `indexes`.
// When the document is not found (nil, nil) is returned.
func MultiGet(ctx context.Context, indexes []Index, id string, dst interface{}, fields ...string) (Index, error) {
	for _, i := range indexes {
		found, err := i.Get(ctx, id, dst, fields...)

		if err != nil {
			return nil, err
		}

		if found {
			return i, nil
		}
	}

	return nil, nil
}
