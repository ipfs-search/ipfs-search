package handler

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/ipfs-search/ipfs-search/components/sniffer/eventsource"

	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
)

// Handler handles EvtProviderPut events by writing Provider's to a channel.
type Handler struct {
	providers chan<- t.Provider
	*instr.Instrumentation
}

// New returns a new handler, writing Provider's to providers.
func New(providers chan<- t.Provider) Handler {
	return Handler{
		providers:       providers,
		Instrumentation: instr.New(),
	}
}

// HandleFunc writes a Provider to the Handler's providers channel for every EvtProviderPut it is called with.
func (h *Handler) HandleFunc(ctx context.Context, e eventsource.EvtProviderPut) error {
	ctx = trace.ContextWithRemoteSpanContext(ctx, e.SpanContext)
	ctx, span := h.Tracer.Start(ctx, "handler.HandleFunc", trace.WithAttributes(
		attribute.Stringer("cid", e.CID),
		attribute.Stringer("peerid", e.PeerID),
	), trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	p := t.Provider{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       e.CID.String(),
		},
		Date:        time.Now(),
		Provider:    e.PeerID.String(),
		SpanContext: span.SpanContext(),
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case h.providers <- p:
		return nil
	}
}
