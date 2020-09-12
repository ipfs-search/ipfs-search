package handler

import (
	"context"
	"log"
	"time"

	t "github.com/ipfs-search/ipfs-search/types"
	"github.com/ipfs-search/ipfs-sniffer/eventsource"
)

type Handler struct {
	providers chan<- t.Provider
}

func New(providers chan<- t.Provider) Handler {
	return Handler{
		providers: providers,
	}
}

func (h *Handler) HandleFunc(ctx context.Context, e eventsource.EvtProviderPut) error {
	log.Printf("Handling event %+v", e)

	p := t.Provider{
		Resource: &t.Resource{
			Protocol: "ipfs",
			ID:       e.CID.String(),
		},
		Date:     time.Now(),
		Provider: e.PeerID.String(),
	}

	log.Printf("Writing provider %+v", p)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case h.providers <- p:
		log.Printf("Written provider %+v", p)
		return nil
	}
}
