package crawler

import (
	"github.com/ipfs-search/ipfs-search/queue"
)

type Queues struct {
	Hashes      queue.Publisher
	Files       queue.Publisher
	Directories queue.Publisher
}
