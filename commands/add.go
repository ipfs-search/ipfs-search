package commands

import (
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/queue"
)

// AddHash queues a single IPFS hash for indexing
func AddHash(hash string) error {
	config, err := getConfig()
	if err != nil {
		return err
	}

	conn, err := queue.NewConnection(config.Factory.AMQPURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.NewChannel()
	if err != nil {
		return err
	}

	queue, err := ch.NewQueue("hashes")
	if err != nil {
		return err
	}

	err = queue.Publish(&crawler.Args{
		Hash: hash,
	})

	return err
}
