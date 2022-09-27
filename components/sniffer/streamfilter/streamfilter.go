package streamfilter

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"

	filters "github.com/ipfs-search/ipfs-search/components/sniffer/providerfilters"
)

// Filter filters a stream of Providers through filters.Filter.
type Filter struct {
	f   filters.Filter
	in  <-chan t.Provider
	out chan<- t.Provider
	*instr.Instrumentation
}

// New creates a new Filter based on a Filter, an incoming and an outgoing channel.
func New(f filters.Filter, in <-chan t.Provider, out chan<- t.Provider) Filter {
	return Filter{
		f:               f,
		in:              in,
		out:             out,
		Instrumentation: instr.New(),
	}
}

func (f *Filter) iterate(ctx context.Context) error {
	select {
	case <-ctx.Done():
		// Context closed, return context error
		return ctx.Err()
	case p := <-f.in:
		return func() error {
			ctx = trace.ContextWithRemoteSpanContext(ctx, p.SpanContext)
			ctx, span := f.Tracer.Start(ctx, "providerfilter.Filter", trace.WithAttributes(
				attribute.String("cid", p.ID),
				attribute.String("peerid", p.Provider),
			))
			defer span.End()

			include, err := f.f.Filter(p)

			if err != nil {
				span.RecordError(err)
			}

			if include {
				// Set new span context on resource
				p.SpanContext = span.SpanContext()

				// Send or return on context close
				select {
				case <-ctx.Done():
					return ctx.Err()
				case f.out <- p:
				}
			}

			return nil
		}()
	}
}

// Filter filters a stream of providers, dropping those for which filter returns false
func (f *Filter) Filter(ctx context.Context) error {
	for {
		if err := f.iterate(ctx); err != nil {
			return err
		}
	}
}
