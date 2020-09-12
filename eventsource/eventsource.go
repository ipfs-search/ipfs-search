package eventsource

import (
	"context"
	"fmt"
	"log"

	"github.com/ipfs-search/ipfs-sniffer/proxy"

	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p-core/event"
)

type EventSource struct {
	bus     event.Bus
	emitter event.Emitter
	ds      datastore.Batching
}

func New(b event.Bus, ds datastore.Batching) (EventSource, error) {
	e, err := b.Emitter(new(EvtProviderPut))
	if err != nil {
		return EventSource{}, err
	}

	s := EventSource{
		bus:     b,
		emitter: e,
	}

	s.ds = proxy.New(ds, s.afterPut)

	return s, nil
}

// nonFatalError is called on non-fatal errors
func (s *EventSource) nonFatalError(err error) {
	log.Printf("error: %v\n", err)
}

func (s *EventSource) afterPut(k datastore.Key, v []byte, err error) error {
	// Ignore error'ed Put's
	if err != nil {
		return err
	}

	// Ignore non-provider keys
	if !isProviderKey(k) {
		return nil
	}

	cid, err := keyToCID(k)
	if err != nil {
		s.nonFatalError(fmt.Errorf("cid from key '%s': %w", k, err))
		return nil
	}

	pid, err := keyToPeerID(k)
	if err != nil {
		s.nonFatalError(fmt.Errorf("pid from key '%s': %w", k, err))
		return nil
	}

	evt := EvtProviderPut{
		CID:    cid,
		PeerID: pid,
	}
	log.Printf("Emitting event: %+v", evt)

	err = s.emitter.Emit(evt)
	if err != nil {
		s.nonFatalError(fmt.Errorf("cid from key '%s': %w", k, err))
		return nil
	}

	return nil
}

func (s *EventSource) Batching() datastore.Batching {
	return s.ds
}

// Subscribe handleFunc to EvtProviderPut events
func (s *EventSource) Subscribe(ctx context.Context, handleFunc func(context.Context, EvtProviderPut) error) error {
	sub, err := s.bus.Subscribe(new(EvtProviderPut))
	if err != nil {
		return fmt.Errorf("subscribing: %w", err)
	}
	defer sub.Close()

	c := sub.Out()
	for {
		log.Println("Waiting for next event")

		select {
		case <-ctx.Done():
			return err
		case e, ok := <-c:
			if !ok {
				return fmt.Errorf("reading from event bus")
			}

			evt, ok := e.(EvtProviderPut)
			if !ok {
				return fmt.Errorf("casting event: %v", evt)
			}

			err := handleFunc(ctx, evt)
			if err != nil {
				return err
			}
		}
	}
}
