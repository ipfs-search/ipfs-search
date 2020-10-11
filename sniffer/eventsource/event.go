package eventsource

import (
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
	"go.opentelemetry.io/otel/api/trace"
)

// EvtProviderPut should be emitted on every datastore Put() for a peer providing a CID.
type EvtProviderPut struct {
	CID         cid.Cid
	PeerID      peer.ID
	SpanContext trace.SpanContext // SpanContext allows a Resource' processing to be traceable across the program
}
