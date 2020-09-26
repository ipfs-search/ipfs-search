package eventsource

// func TestProviderKeyToCIDNamespacesError(t *testing.T) {
// 	_, err := providerKeyToCID(datastore.NewKey("invalid"))
// 	if err != errInvalidKeyNamespaces {
// 		t.Fatal("expected invalid key namespaces error")
// 	}
// }

// func TestProviderKeyToCIDEncodingBase32Error(t *testing.T) {
// 	_, err := providerKeyToCID(datastore.NewKey("/providers/8"))
// 	if err == nil {
// 		t.Fatal("expected invalid base32 encoding error")
// 	}
// }

// func TestProviderKeyToCIDEncodingCIDError(t *testing.T) {
// 	_, err := providerKeyToCID(datastore.NewKey("/providers/base32notcid"))
// 	if err == nil {
// 		t.Fatal("expected invalid CID encoding error")
// 	}
// }
