package handler

import (
	"context"
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
	p := t.Provider{
		Resource: &t.Resource{
			Protocol: "ipfs",
			ID:       e.CID.String(),
		},
		Date:     time.Now(),
		Provider: e.PeerID.String(),
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case h.providers <- p:
		return nil
	}
}
