package pool

import (
	"net/http"

	"github.com/ipfs-search/ipfs-search/components/protocol"
	"github.com/ipfs-search/ipfs-search/components/protocol/ipfs"
	"github.com/ipfs-search/ipfs-search/utils"
)

func (p *Pool) getProtocol() protocol.Protocol {
	ipfsTransport := utils.GetHTTPTransport(p.dialer.DialContext, p.config.Workers.MaxIPFSConns)
	ipfsClient := &http.Client{Transport: ipfsTransport}

	return ipfs.New(p.config.IPFSConfig(), ipfsClient, p.Instrumentation)
}
