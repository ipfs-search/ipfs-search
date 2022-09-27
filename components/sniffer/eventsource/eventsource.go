package eventsource

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-eventbus"
	"github.com/libp2p/go-libp2p-core/event"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/ipfs-search/ipfs-search/components/sniffer/proxy"
	"github.com/ipfs-search/ipfs-search/instr"
)

const bufSize = 512

var handleTimeout = time.Second

type handleFunc func(context.Context, EvtProviderPut) error

// EventSource generates events on a Bus after Put operations by proxying a Batching datastore.
type EventSource struct {
	bus     event.Bus
	emitter event.Emitter
	ds      datastore.Batching
	*instr.Instrumentation
}

// New sets up a new EventSource or returns an error.
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
	_, span := s.Tracer.Start(context.TODO(), "eventsource.afterPut")
	defer span.End()

	// Ignore error'ed Put's
	if err != nil {
		span.RecordError(err)
		return err
	}

	// Ignore non-provider keys
	if !isProviderKey(k) {
		span.RecordError(fmt.Errorf("Non-provider key"))
		return nil
	}

	cid, err := keyToCID(k)
	if err != nil {
		span.RecordError(fmt.Errorf("cid from key '%s': %w", k, err))
		return nil
	}

	pid, err := keyToPeerID(k)
	if err != nil {
		span.RecordError(fmt.Errorf("pid from key '%s': %w", k, err))
		return nil
	}

	span.SetAttributes(
		attribute.Stringer("cid", cid),
		attribute.Stringer("peerid", pid),
	)

	e := EvtProviderPut{
		CID:         cid,
		PeerID:      pid,
		SpanContext: span.SpanContext(),
	}

	if err := s.emitter.Emit(e); err != nil {
		span.RecordError(err)
	} else {
		span.SetStatus(codes.Ok, "emitted")
	}

	// Return *original* error
	return err
}

// Batching returns the proxied Batching datastore.
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

// Subscribe handleFunc to EvtProviderPut events.
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
