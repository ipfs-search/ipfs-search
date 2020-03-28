package sniffer

import (
	"github.com/ipfs/go-cid"
	"log"
)

type cidFilter struct {
}

func NewCidFilter() *cidFilter {
	return &cidFilter{}
}

func (f *cidFilter) Filter(p Provider) bool {
	c, err := cid.Decode(p.Id)

	if err != nil {
		// Error: resource is not a valid CID
		log.Printf("Error decoding %v to CID: %v", p, err)
		return false
	}

	switch t := c.Type(); t {
	case cid.Raw, cid.DagProtobuf:
		// (Potential) files and directories
		return true
	default:
		// Can't handle other types (for now)
		log.Printf("Unsupported codec %v, codec %v", p, cid.CodecToStr[t])
		return false
	}
}
