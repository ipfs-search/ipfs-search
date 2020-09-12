package sniffer

import (
	"context"
	"fmt"
	"log"

	"github.com/ipfs-search/ipfs-sniffer/eventsource"
	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-eventbus"
)

type Sniffer struct {
	es eventsource.EventSource
}

func New(ds datastore.Batching) (*Sniffer, error) {
	// Test update
	bus := eventbus.NewBus()

	es, err := eventsource.New(bus, ds)
	if err != nil {
		return nil, fmt.Errorf("failed to get eventsource: %w", err)
	}

	s := Sniffer{
		es: es,
	}

	return &s, nil
}

func (s *Sniffer) Batching() datastore.Batching {
	return s.es.Batching()
}

func (s *Sniffer) Sniff(ctx context.Context) error {
	eventHandler := func(ctx context.Context, e eventsource.EvtProviderPut) error {
		log.Printf("%+v\n", e)
		return nil
	}

	return s.es.Subscribe(ctx, eventHandler)
}
