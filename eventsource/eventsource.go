package eventsource

import (
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

	err = s.emitter.Emit(EvtProviderPut{
		CID:    cid,
		PeerID: pid,
	})
	if err != nil {
		s.nonFatalError(fmt.Errorf("cid from key '%s': %w", k, err))
		return nil
	}

	return nil
}

func (s *EventSource) Batching() datastore.Batching {
	return s.ds
}

// Subscribe to EvtProviderPut events. Don't forget to call Close() on the subscription!
func (s *EventSource) Subscribe() (event.Subscription, error) {
	return s.bus.Subscribe(new(EvtProviderPut))
}
