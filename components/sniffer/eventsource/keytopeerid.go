package eventsource

import (
	"strings"

	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-base32"
)

// Source: https://github.com/libp2p/go-libp2p-kad-dht/blob/9304f5575ea4c578d1316c2cf695a06b65c88dbe/providers/providers_manager.go#L339
func keyToPeerID(k datastore.Key) (peer.ID, error) {
	kStr := k.String()

	lix := strings.LastIndex(kStr, "/")

	decstr, err := base32.RawStdEncoding.DecodeString(kStr[lix+1:])
	if err != nil {
		return "", err
	}

	return peer.ID(decstr), nil
}
