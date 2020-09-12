package sniffer

import (
	"context"
	"fmt"

	"github.com/ipfs-search/ipfs-sniffer/eventsource"
)

type Sniffer struct {
	es *eventsource.EventSource
}

func New(es *eventsource.EventSource) Sniffer {
	return Sniffer{
		es: es,
	}
}

func (s *Sniffer) Sniff(ctx context.Context) error {
	sub, err := s.es.Subscribe()
	if err != nil {
		return fmt.Errorf("subscribing: %w", err)
	}
	defer sub.Close()

	c := sub.Out()
	for {
		select {
		case <-ctx.Done():
			return err
		case e, ok := <-c:
			if !ok {
				return fmt.Errorf("reading from event bus")
			}

			evt, ok := e.(eventsource.EvtProviderPut)
			if !ok {
				return fmt.Errorf("casting event: %v", evt)
			}
		}
	}
}
