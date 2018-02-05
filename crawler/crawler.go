package crawler

import (
	"encoding/json"
	"github.com/ipfs-search/ipfs-search/indexer"
	"github.com/ipfs-search/ipfs-search/queue"
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

	Shell     *shell.Shell
	Indexer   *indexer.Indexer
	FileQueue *queue.Queue
	HashQueue *queue.Queue
}

// IndexableFromJSON returns and Indexable associated with this crawler based on a JSON blob
func (c *Crawler) IndexableFromJSON(input []byte) (*Indexable, error) {
	// Unmarshall message into crawler Args
	args := &Args{}
	err := json.Unmarshal(input, args)
	if err != nil {
		return nil, err
	}

	return &Indexable{
		Args:    args,
		Crawler: c,
	}, nil

}
