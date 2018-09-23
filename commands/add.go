package commands

import (
	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/queue"
)

// AddHash queues a single IPFS hash for indexing
func AddHash(hash string) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	conn, err := queue.NewConnection(cfg.AMPQ.AMQPURL)
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
