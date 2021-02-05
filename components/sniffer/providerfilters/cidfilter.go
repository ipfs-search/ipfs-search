package providerfilters

import (
	"errors"
	"fmt"

	"github.com/ipfs/go-cid"

	t "github.com/ipfs-search/ipfs-search/types"
)

var (
	errUnsupportedProtocol = errors.New("unsupported protocol")
	errDecodingCID         = errors.New("unable to decode CID")
	errUnsupportedCodec    = errors.New("unsupported codec")
)

// CidFilter filters out invalid CID's or those which are not Raw or DagProtobuf.
type CidFilter struct{}

// NewCidFilter returns a pointer to a new CidFilter.
func NewCidFilter() *CidFilter {
	return &CidFilter{}
}

// Filter takes a Provider and returns true when it is to be included, false
// when not and an error when unexpected condition occur.
func (f *CidFilter) Filter(p t.Provider) (bool, error) {
	if p.Resource.Protocol != t.IPFSProtocol {
		return false, fmt.Errorf("%w: %s for %v", errUnsupportedProtocol, p.Resource.Protocol, p)
	}

	c, err := cid.Decode(p.ID)

	if err != nil {
		return false, fmt.Errorf("%w: %s decoding CID %v", errDecodingCID, err, p)
	}

	switch cidType := c.Type(); cidType {
	case cid.Raw, cid.DagProtobuf:
		// (Potential) files and directories
		return true, nil
	default:
		// Can't handle other types (for now)
		return false, fmt.Errorf("%w: %s for %v", errUnsupportedCodec, cid.CodecToStr[cidType], p)
	}
}
