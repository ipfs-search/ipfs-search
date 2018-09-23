package commands

import (
	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/queue"
)

// AddHash queues a single IPFS hash for indexing
func AddHash(cfg *config.Config, hash string) error {
	conn, err := queue.NewConnection(cfg.AMQP.AMQPURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	queue, err := conn.NewChannelQueue("hashes")
	if err != nil {
		return err
	}

	err = queue.Publish(&crawler.Args{
		Hash: hash,
	})

	return err
}
