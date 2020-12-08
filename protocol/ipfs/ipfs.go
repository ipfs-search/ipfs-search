package ipfs

import (
	"fmt"
	"net/http"
	"net/url"

	ipfs "github.com/ipfs/go-ipfs-api"

	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs-search/ipfs-search/protocol"
	t "github.com/ipfs-search/ipfs-search/types"
)

// IPFS implements the Protocol interface for the Interplanery Filesystem. It is concurrency-safe.
type IPFS struct {
	config *Config

	gatewayURL *url.URL
	shell      *ipfs.Shell

	*instr.Instrumentation
}

// absolutePath returns the absolute (CID-only) path for a resource.
func absolutePath(r *t.AnnotatedResource) string {
	return fmt.Sprintf("/ipfs/%s", r.ID)
}

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

// New returns a new IPFS protocol.
func New(config *Config, client *http.Client, instr *instr.Instrumentation) *IPFS {
	// Initialize gatewayURL
	gatewayURL, err := url.Parse(config.IPFSGatewayURL)
	if err != nil {
		panic(fmt.Sprintf("could not parse IPFS Gateway URL, error: %v", err))
	}

	if !gatewayURL.IsAbs() {
		panic(fmt.Sprintf("gateway URL is not absolute: %s", gatewayURL))
	}

	// Create IPFS shell
	shell := ipfs.NewShellWithClient(config.IPFSAPIURL, client)

	return &IPFS{
		config,
		gatewayURL,
		shell,
		instr,
	}
}

// Compile-time assurance that implementation satisfies interface.
var _ protocol.Protocol = &IPFS{}
