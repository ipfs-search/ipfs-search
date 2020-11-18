package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/ipfs/go-ipfs-api"

	"github.com/ipfs-search/ipfs-search/extractor"
	index_types "github.com/ipfs-search/ipfs-search/index/types"
	"github.com/ipfs-search/ipfs-search/protocol"
	t "github.com/ipfs-search/ipfs-search/types"
)

type Crawler struct {
	indexes   Indexes
	protocol  protocol.Protocol
	extractor extractor.Extractor
}

// Crawler consumes file and hash queues and indexes them
type Crawler struct {
	Config *Config

	Shell     *shell.Shell
	Extractor extractor.Extractor

	FileIndex      index.Index
	DirectoryIndex index.Index
	InvalidIndex   index.Index

	FileQueue queue.Publisher
	HashQueue queue.Publisher

	*instr.Instrumentation
}

// IndexableFromJSON returns and Indexable associated with this crawler based on a JSON blob
func (c *Crawler) IndexableFromJSON(input []byte) (*Indexable, error) {
	// Unmarshall message into crawler Args
	args := &Args{}
	err := json.Unmarshal(input, args)
	if err != nil {
		return nil, err
	}

	// Later down, we assume this hash is set and we're seeing errors where
	// this aparently seems not the case.
	if args.Hash == "" {
		return nil, fmt.Errorf("Empty hash in JSON: %s", input)
	}

	return &Indexable{
		Args:    args,
		Crawler: c,
	}, nil

}
