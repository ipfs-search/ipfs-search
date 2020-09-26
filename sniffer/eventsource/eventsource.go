package eventsource

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ipfs-search/ipfs-sniffer/proxy"

	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-eventbus"
	"github.com/libp2p/go-libp2p-core/event"
)

const bufSize = 256

var handleTimeout = time.Second

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

	e := EvtProviderPut{
		CID:    cid,
		PeerID: pid,
	}

	log.Printf("Emitting Put Event %s", e)

	err = s.emitter.Emit(e)
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
// TODO: Make this return errgroup, err instead of blocking - leaving the caller to decide how to deal with it and separating
// initialisation from listening.
func (s *EventSource) Subscribe(ctx context.Context, handleFunc func(context.Context, EvtProviderPut) error) error {
	sub, err := s.bus.Subscribe(new(EvtProviderPut), eventbus.BufSize(bufSize))
	if err != nil {
		return fmt.Errorf("subscribing: %w", err)
	}
	defer sub.Close()

	c := sub.Out()

	// TODO: Consider running this in a Goroutine through an errorgroup
	for {
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

			// Timeout handler to expose issues on the handler side
			ctx, cancel := context.WithTimeout(ctx, handleTimeout)

			err := handleFunc(ctx, evt)
			cancel() // Frees up timeout context's resources

			if err != nil {
				return err
			}
		}
	}
}
