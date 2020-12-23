package crawler

import (
	"github.com/ipfs-search/ipfs-search/queue"
)

// Queues used for crawling.
type Queues struct {
	Files       queue.Queue
	Directories queue.Queue
	Hashes      queue.Queue
}
