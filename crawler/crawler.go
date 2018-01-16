package crawler

import (
	"github.com/ipfs-search/ipfs-search/indexer"
	"github.com/ipfs-search/ipfs-search/queue"
	"github.com/ipfs/go-ipfs-api"
	"log"
	"net"
	"net/url"
	// "path"
	"strings"
	"syscall"
	"time"
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

// Handle IPFS errors graceously, returns try again bool and original error
func (c *Crawler) handleError(err error, hash string) (bool, error) {
	if _, ok := err.(*shell.Error); ok && strings.Contains(err.Error(), "proto") {
		// We're not recovering from protocol errors, so panic

		// Attempt to index panic to prevent re-indexing
		m := metadata{
			"error": err.Error(),
		}

		c.Indexer.IndexItem("invalid", hash, m)

		panic(err)
	}

	if uerr, ok := err.(*url.Error); ok {
		// URL errors

		log.Printf("URL error: %v", uerr)

		if uerr.Timeout() {
			// Fail on timeouts
			return false, err
		}

		if uerr.Temporary() {
			// Retry on other temp errors
			return true, nil
		}

		// Somehow, the errors below are not temp errors !?
		switch t := uerr.Err.(type) {
		case *net.OpError:
			if t.Op == "dial" {
				log.Printf("Unknown host %v", t)
				return true, nil

			} else if t.Op == "read" {
				log.Printf("Connection refused %v", t)
				return true, nil
			}

		case syscall.Errno:
			if t == syscall.ECONNREFUSED {
				log.Printf("Connection refused %v", t)
				return true, nil
			}
		}
	}

	return false, err
}

// CrawlHash crawls a particular hash (file or directory)
func (c *Crawler) CrawlHash(args *Args) error {
	references, alreadyIndexed, err := c.indexReferences(args.Hash, args.Name, args.ParentHash)

	if err != nil {
		return err
	}

	if alreadyIndexed {
		return nil
	}

	log.Printf("Crawling hash '%s' (%s)", args.Hash, args.Name)

	url := hashURL(args.Hash)

	var list *shell.UnixLsObject

	tryAgain := true
	for tryAgain {
		list, err = c.Shell.FileList(url)

		tryAgain, err = c.handleError(err, args.Hash)

		if tryAgain {
			log.Printf("Retrying in %s", c.Config.RetryWait)
			time.Sleep(c.Config.RetryWait)
		}
	}

	if err != nil {
		return err
	}

	switch list.Type {
	case "File":
		// Add to file crawl queue
		fileArgs := Args{
			Hash:       args.Hash,
			Name:       args.Name,
			Size:       list.Size,
			ParentHash: args.ParentHash,
		}

		err = c.FileQueue.AddTask(fileArgs)
		if err != nil {
			// failed to send the task
			return err
		}
	case "Directory":
		// Queue indexing of linked items
		for _, link := range list.Links {
			dirArgs := Args{
				Hash:       link.Hash,
				Name:       link.Name,
				Size:       link.Size,
				ParentHash: args.Hash,
			}

			switch link.Type {
			case "File":
				// Add file to crawl queue
				err = c.FileQueue.AddTask(dirArgs)
				if err != nil {
					// failed to send the task
					return err
				}

			case "Directory":
				// Add directory to crawl queue
				c.HashQueue.AddTask(dirArgs)
				if err != nil {
					// failed to send the task
					return err
				}
			default:
				log.Printf("Type '%s' skipped for '%s'", link.Type, args.Hash)
			}
		}

		// Index name and size for directory and directory items
		properties := map[string]interface{}{
			"links":      list.Links,
			"size":       list.Size,
			"references": references,
		}

		// Skip partial content
		if list.Size == c.Config.PartialSize && args.ParentHash == "" {
			// Assertion error.
			// REMOVE ME!
			log.Printf("Skipping unreferenced partial content for directory %s", args.Hash)
			return nil
		}

		err := c.Indexer.IndexItem("directory", args.Hash, properties)
		if err != nil {
			return err
		}

	default:
		log.Printf("Type '%s' skipped for '%s'", list.Type, args.Hash)
	}

	log.Printf("Finished hash %s", args.Hash)

	return nil
}

// CrawlFile crawls a single object, known to be a file
func (c *Crawler) CrawlFile(args *Args) error {
	if args.Size == c.Config.PartialSize && args.ParentHash == "" {
		// Assertion error.
		// REMOVE ME!
		log.Printf("Skipping unreferenced partial content for file %s", args.Hash)
		return nil
	}

	references, alreadyIndexed, err := c.indexReferences(args.Hash, args.Name, args.ParentHash)

	if err != nil {
		return err
	}

	if alreadyIndexed {
		return nil
	}

	log.Printf("Crawling file %s (%s)\n", args.Hash, args.Name)

	m := make(metadata)
	c.getMetadata(args, &m)
	if err != nil {
		return err
	}

	// Add previously found references now
	m["size"] = args.Size
	m["references"] = references

	err = c.Indexer.IndexItem("file", args.Hash, m)
	if err != nil {
		return err
	}

	log.Printf("Finished file %s", args.Hash)

	return nil
}
