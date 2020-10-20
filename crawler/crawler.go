package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/ipfs-search/ipfs-search/index"
	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs/go-ipfs-api"
)

// Args describe a resource to be crawled
type Args struct {
	Hash       string
	Name       string
	Size       uint64
	ParentHash string
	ParentName string // This is legacy, should be removed
}

// Crawler consumes file and hash queues and indexes them
type Crawler struct {
	Config *Config

	Shell *shell.Shell

	FileIndex      index.Index
	DirectoryIndex index.Index
	InvalidIndex   index.Index

	FileQueue Queue
	HashQueue Queue

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
