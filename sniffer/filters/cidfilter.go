package filters

import (
	"fmt"
	t "github.com/ipfs-search/ipfs-search/types"
	"github.com/ipfs/go-cid"
)

type cidFilter struct{}

func CidFilter() *cidFilter {
	return &cidFilter{}
}

func (f *cidFilter) Filter(p t.Provider) (bool, error) {
	c, err := cid.Decode(p.Id)

	if err != nil {
		return false, fmt.Errorf("Error decoding %v to CID: %v", p, err)
	}

	switch t := c.Type(); t {
	case cid.Raw, cid.DagProtobuf:
		// (Potential) files and directories
		return true, nil
	default:
		// Can't handle other types (for now)
		return false, fmt.Errorf("Unsupported codec %v, codec %v", p, cid.CodecToStr[t])
	}
}
