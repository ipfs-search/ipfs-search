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

func New(ds datastore.Batching) (*Sniffer, datastore.Batching, error) {
	bus := eventbus.NewBus()

	es, err := eventsource.New(bus, ds)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get eventsource: %w", err)
	}

	s := Sniffer{
		es: es,
	}

	return &s, es.Batching(), nil
}

func (s Sniffer) Sniff(ctx context.Context) error {
	eventHandler := func(ctx context.Context, e eventsource.EvtProviderPut) error {
		log.Printf("%+v\n", e)
		return nil
	}

	return s.es.Subscribe(ctx, eventHandler)
}
