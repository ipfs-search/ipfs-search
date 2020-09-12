package sniffer

import (
	"context"
	"log"

	"github.com/ipfs-search/ipfs-sniffer/eventsource"
)

type Sniffer struct {
	es eventsource.EventSource
}

func New(es eventsource.EventSource) Sniffer {
	return Sniffer{
		es: es,
	}
}

func (s Sniffer) Sniff(ctx context.Context) error {
	eventHandler := func(ctx context.Context, e eventsource.EvtProviderPut) error {
		log.Printf("%v\n", e)
		return nil
	}

	return s.es.Subscribe(ctx, eventHandler)
}
