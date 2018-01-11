package commands

import (
	"github.com/ipfs-search/ipfs-search/queue"
)

// AddHash queues a single IPFS hash for indexing
func AddHash(hash string) error {
	ch, err := queue.NewChannel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queue, err := queue.NewTaskQueue(ch, "hashes")
	if err != nil {
		return err
	}

	err = queue.AddTask(map[string]interface{}{
		"hash": hash,
	})

	return err
}
