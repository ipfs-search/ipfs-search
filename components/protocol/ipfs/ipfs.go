package ipfs

import (
	"fmt"
	"net/http"
	"net/url"

	ipfs "github.com/ipfs/go-ipfs-api"

	"github.com/ipfs-search/ipfs-search/components/protocol"

	"github.com/ipfs-search/ipfs-search/instr"
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

// New returns a new IPFS protocol.
func New(config *Config, client *http.Client, instr *instr.Instrumentation) *IPFS {
	// Initialize gatewayURL
	gatewayURL, err := url.Parse(config.GatewayURL)
	if err != nil {
		panic(fmt.Sprintf("could not parse IPFS Gateway URL, error: %v", err))
	}

	if !gatewayURL.IsAbs() {
		panic(fmt.Sprintf("gateway URL is not absolute: %s", gatewayURL))
	}

	// Create IPFS shell
	shell := ipfs.NewShellWithClient(config.APIURL, client)

	return &IPFS{
		config,
		gatewayURL,
		shell,
		instr,
	}
}

// Compile-time assurance that implementation satisfies interface.
var _ protocol.Protocol = &IPFS{}
