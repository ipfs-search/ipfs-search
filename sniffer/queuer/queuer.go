package queuer

import (
	"context"
	"time"

	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs-search/ipfs-search/queue"
	t "github.com/ipfs-search/ipfs-search/types"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
)

type Queuer struct {
	queue        queue.Publisher
	providers    <-chan t.Provider
	queueTimeout time.Duration
	*instr.Instrumentation
}

func New(q queue.Publisher, providers <-chan t.Provider) Queuer {
	return Queuer{
		queue:           q,
		providers:       providers,
		queueTimeout:    5 * time.Minute, // Kamikaze after 5 minutes of waiting
		Instrumentation: instr.New(),
	}
}

func (q *Queuer) iterate(ctx context.Context) error {
	// Never wait more than queueTimeout for a message
	ctx, cancel := context.WithTimeout(ctx, q.queueTimeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case p := <-q.providers:
		return func() error {
			ctx = trace.ContextWithRemoteSpanContext(ctx, p.SpanContext)
			_, span := q.Tracer.Start(ctx, "queue.Publish", trace.WithAttributes(
				label.String("cid", p.ID),
				label.String("peerid", p.Provider),
			), trace.WithSpanKind(trace.SpanKindProducer))
			defer span.End()

			// TODO: Provider channel should be pointer stream, preventing copying of data.

			// Add with highest priority (9), as this is supposed to be available
			err := q.queue.Publish(ctx, &p, 9)

			if err != nil {
				span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
			} else {
				span.SetStatus(codes.Ok, "published")
			}

			return err
		}()
	}
}

func (q *Queuer) Queue(ctx context.Context) error {
	for {
		if err := q.iterate(ctx); err != nil {
			return err
		}
	}
}
