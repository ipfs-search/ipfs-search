package ipfs

import (
	"fmt"
	"net/url"

	t "github.com/ipfs-search/ipfs-search/types"
)

// namedPath returns the (escaped/raw) path for a resource.
// If a reference is available, it is used to generate the filename to facilitate content
// type detection (e.g. /ipfs/<parent_hash>/my_file.jpg instead of /ipfs/<file_hash>/).
func namedPath(r *t.AnnotatedResource) string {
	if ref := r.Reference; ref.Name != "" {
		return fmt.Sprintf("/ipfs/%s/%s", ref.Parent.ID, url.PathEscape(ref.Name))
	}

	return absolutePath(r)
}

// GatewayURL returns the URL to request a resource from the gateway.
// If a reference is available, it is used to generate the filename to facilitate content
// type detection (e.g. /ipfs/<parent_hash>/my_file.jpg instead of /ipfs/<file_hash>/).
// Ref: http://docs.ipfs.io.ipns.localhost:8080/concepts/ipfs-gateway/#gateway-types
func (i *IPFS) GatewayURL(r *t.AnnotatedResource) string {
	url, err := i.gatewayURL.Parse(namedPath(r))

	if err != nil {
		panic(fmt.Sprintf("error generating GatewayURL: %v", err))
	}

	return url.String()
}
