package eventsource

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-eventbus"
	"github.com/libp2p/go-libp2p-core/event"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"

	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs-search/ipfs-search/sniffer/proxy"
)

const bufSize = 512

var handleTimeout = time.Second

type handleFunc func(context.Context, EvtProviderPut) error

type EventSource struct {
	bus     event.Bus
	emitter event.Emitter
	ds      datastore.Batching
	*instr.Instrumentation
}

func New(b event.Bus, ds datastore.Batching) (EventSource, error) {
	e, err := b.Emitter(new(EvtProviderPut))
	if err != nil {
		return EventSource{}, err
	}

	s := EventSource{
		bus:             b,
		emitter:         e,
		Instrumentation: instr.New(),
	}

	s.ds = proxy.New(ds, s.afterPut)

	return s, nil
}

func (s *EventSource) afterPut(k datastore.Key, v []byte, err error) error {
	ctx, span := s.Tracer.Start(context.TODO(), "eventsource.afterPut")
	defer span.End()

	// Ignore error'ed Put's
	if err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Ok))
		return err
	}

	// Ignore non-provider keys
	if !isProviderKey(k) {
		span.RecordError(ctx, fmt.Errorf("Non-provider key"), trace.WithErrorStatus(codes.Ok))
		return nil
	}

	cid, err := keyToCID(k)
	if err != nil {
		span.RecordError(ctx, fmt.Errorf("cid from key '%s': %w", k, err), trace.WithErrorStatus(codes.Error))
		return nil
	}

	pid, err := keyToPeerID(k)
	if err != nil {
		span.RecordError(ctx, fmt.Errorf("pid from key '%s': %w", k, err), trace.WithErrorStatus(codes.Error))
		return nil
	}

	span.SetAttributes(
		label.Stringer("cid", cid),
		label.Stringer("peerid", pid),
	)

	e := EvtProviderPut{
		CID:         cid,
		PeerID:      pid,
		SpanContext: span.SpanContext(),
	}

	if err := s.emitter.Emit(e); err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
	} else {
		span.SetStatus(codes.Ok, "emitted")
	}

	// Return *original* error
	return err
}

func (s *EventSource) Batching() datastore.Batching {
	return s.ds
}

func (s *EventSource) iterate(ctx context.Context, c <-chan interface{}, h handleFunc) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case e, ok := <-c:
		if !ok {
			return fmt.Errorf("reading from event bus")
		}

		evt, ok := e.(EvtProviderPut)
		if !ok {
			return fmt.Errorf("casting event: %v", evt)
		}

		// Timeout handler to expose issues on the handler side
		ctx, cancel := context.WithTimeout(ctx, handleTimeout)
		defer cancel()

		return h(ctx, evt)
	}
}

// Subscribe handleFunc to EvtProviderPut events
// TODO: Make this return errgroup, err instead of blocking - leaving the caller to decide how to deal with it and separating
// initialisation from listening.
func (s *EventSource) Subscribe(ctx context.Context, h handleFunc) error {
	sub, err := s.bus.Subscribe(new(EvtProviderPut), eventbus.BufSize(bufSize))
	if err != nil {
		return fmt.Errorf("subscribing: %w", err)
	}
	defer sub.Close()

	c := sub.Out()

	for {
		if err := s.iterate(ctx, c, h); err != nil {
			return err
		}
	}
}
