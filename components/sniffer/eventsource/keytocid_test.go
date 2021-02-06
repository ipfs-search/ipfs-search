package eventsource

import (
	"github.com/ipfs/go-datastore"
	"testing"
)

func TestKeyToCIDNamespacesError(t *testing.T) {
	_, err := keyToCID(datastore.NewKey("invalid"))
	if err != errInvalidKeyNamespaces {
		t.Fatal("expected invalid key namespaces error")
	}
}

func TestKeyToCIDEncodingBase32Error(t *testing.T) {
	_, err := keyToCID(datastore.NewKey("/providers/8"))
	if err == nil {
		t.Fatal("expected invalid base32 encoding error")
	}
}

func TestKeyToCIDEncodingCIDError(t *testing.T) {
	_, err := keyToCID(datastore.NewKey("/providers/base32notcid"))
	if err == nil {
		t.Fatal("expected invalid CID encoding error")
	}
}
