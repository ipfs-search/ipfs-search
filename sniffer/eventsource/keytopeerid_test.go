package eventsource

import (
	"github.com/ipfs/go-datastore"
	"testing"
)

func TestKeyToPeerID(t *testing.T) {
	// CID: QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp
	// Provider: QmeTtFXm42Jb2todcKR538j6qHYxXt6suUzpF3rtT9FPSd
	key := datastore.NewKey("/providers/CIQDWKPBHXLJ3XVELRJZA2SYY7OGCSX6FRSIZS2VQQPVKOA2Z4VXN2I/CIQO7FK6IWMEVZU2QU6QRJKMCLW4DXQGSVSVB3V56Y272TB3IPSBGFQ")

	id, err := keyToPeerID(key)
	if err != nil {
		t.Fatal(err)
	}

	if id.String() != "QmeTtFXm42Jb2todcKR538j6qHYxXt6suUzpF3rtT9FPSd" {
		t.Fatal("Peer ID not equal")
	}
}
