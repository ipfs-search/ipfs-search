package handler

import (
	"context"
	"log"
	"time"

	"github.com/ipfs-search/ipfs-search/sniffer/eventsource"
	t "github.com/ipfs-search/ipfs-search/types"
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
	p := t.Provider{
		Resource: &t.Resource{
			Protocol: "ipfs",
			ID:       e.CID.String(),
		},
		Date:     time.Now(),
		Provider: e.PeerID.String(),
	}

	log.Printf("Handling provider %s", p)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case h.providers <- p:
		return nil
	}
}
