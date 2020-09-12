package eventsource

import (
	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p-kad-dht/providers"
)

// root namespace of provider keys
var providersRoot = datastore.NewKey(providers.ProvidersKeyPrefix)

func isProviderKey(k datastore.Key) bool {
	// not interested if this is not a query for providers of a particular cid
	// we're looking for /providers/cid, not /providers (currently used in GC)
	if !providersRoot.IsAncestorOf(k) || len(k.Namespaces()) < 2 {
		return false
	}

	return true
}
