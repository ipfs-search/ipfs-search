package crawler

import (
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
	FileQueue *queue.TaskQueue
	HashQueue *queue.TaskQueue
}
