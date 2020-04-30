package filters

import (
	t "github.com/ipfs-search/ipfs-search/types"
	"github.com/ipfs/go-cid"
)

// CidFilter filters out invalid CID's or those which are not Raw or DagProtobuf.
type CidFilter struct{}

// NewCidFilter returns a new CidFilter.
func NewCidFilter() *CidFilter {
	return &CidFilter{}
}

// Filter filters out invalid CID's or those which are not Raw or DagProtobuf.
func (f *CidFilter) Filter(p t.Provider) (bool, error) {
	if p.Resource.Protocol != "ipfs" {
		return false, t.NewProviderErrorf(nil, p, "Unsupported protocol %s for %v", p.Resource.Protocol, p)
	}

	c, err := cid.Decode(p.Id)

	if err != nil {
		return false, t.NewProviderErrorf(err, p, "%s decoding CID %v", err, p)
	}

	switch cidType := c.Type(); cidType {
	case cid.Raw, cid.DagProtobuf:
		// (Potential) files and directories
		return true, nil
	default:
		// Can't handle other types (for now)
		return false, t.NewProviderErrorf(nil, p, "Unsupported codec %s for %v", cid.CodecToStr[cidType], p)
	}
}
