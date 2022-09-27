package extractor

import (
	"context"
	"fmt"

	"github.com/c2h5oh/datasize"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	t "github.com/ipfs-search/ipfs-search/types"
)

// ValidateMaxSize returns ErrFileTooLarge when the resource size is above maxSize.
func ValidateMaxSize(ctx context.Context, r *t.AnnotatedResource, maxSize datasize.ByteSize) error {
	span := trace.SpanFromContext(ctx)

	if r.Size > uint64(maxSize) {
		err := fmt.Errorf("%w: %d", ErrFileTooLarge, r.Size)
		span.RecordError(ErrFileTooLarge,
			trace.WithAttributes(
				attribute.Int64("file.size", int64(r.Size)),
			),
		)
		return err
	}

	return nil
}
