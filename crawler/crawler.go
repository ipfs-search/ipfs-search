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

// skipItem determines whether a particular item should not be indexed
// This holds particularly to partial content.
func (c *Crawler) skipItem(args *Args) bool {
	if args.Size == c.Config.PartialSize && args.ParentHash == "" {
		log.Printf("Skipping unreferenced partial content for file %s", args.Hash)
		return true
	}

	return false
}

// Handle errors graceously, returns try again bool and original error
// TODO: this handles both errors for listing as well as metadata errors,
// which seems a very bad idea and makes this function unnecessarily complex.
// We should figure out which code handles which and split it up.
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

// getFileList return list of files and/or type of item (directory/file)
func (c *Crawler) getFileList(args *Args) (list *shell.UnixLsObject, err error) {
	url := hashURL(args.Hash)

	tryAgain := true
	for tryAgain {
		list, err = c.Shell.FileList(url)

		tryAgain, err = c.handleError(err, args.Hash)

		if tryAgain {
			log.Printf("Retrying in %s", c.Config.RetryWait)
			time.Sleep(c.Config.RetryWait)
		}
	}

	return
}

// queueList queues any items in a given list/directory
func (c *Crawler) queueList(args *Args, list *shell.UnixLsObject) (err error) {
	for _, link := range list.Links {
		dirArgs := &Args{
			Hash:       link.Hash,
			Name:       link.Name,
			Size:       link.Size,
			ParentHash: args.Hash,
		}

		switch link.Type {
		case "File":
			// Add file to crawl queue
			err = c.FileQueue.AddTask(dirArgs)
		case "Directory":
			// Add directory to crawl queue
			err = c.HashQueue.AddTask(dirArgs)
		default:
			log.Printf("Type '%s' skipped for '%s'", link.Type, args.Hash)
		}
	}

	return
}

// processList processes a file listing
func (c *Crawler) processList(args *Args, list *shell.UnixLsObject, references []indexer.Reference) (err error) {
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
	case "Directory":
		// Queue indexing of linked items
		err = c.queueList(args, list)
		if err != nil {
			return err
		}

		// Index name and size for directory and directory items
		properties := metadata{
			"links":      list.Links,
			"size":       list.Size,
			"references": references,
		}

		err = c.Indexer.IndexItem("directory", args.Hash, properties)
	default:
		log.Printf("Type '%s' skipped for '%s'", list.Type, args.Hash)
	}

	return
}

// CrawlHash crawls a particular hash (file or directory)
func (c *Crawler) CrawlHash(args *Args) error {
	if c.skipItem(args) {
		return nil
	}

	references, alreadyIndexed, err := c.indexReferences(args.Hash, args.Name, args.ParentHash)
	if err != nil {
		return err
	}

	if alreadyIndexed {
		return nil
	}

	log.Printf("crawling hash '%s' (%s)", args.Hash, args.Name)

	list, err := c.getFileList(args)
	if err != nil {
		return err
	}

	err = c.processList(args, list, references)
	if err != nil {
		return err
	}

	log.Printf("Finished hash %s", args.Hash)

	return nil
}

// CrawlFile crawls a single object, known to be a file
func (c *Crawler) CrawlFile(args *Args) error {
	if c.skipItem(args) {
		return nil
	}

	references, alreadyIndexed, err := c.indexReferences(args.Hash, args.Name, args.ParentHash)

	if err != nil {
		return err
	}

	if alreadyIndexed {
		return nil
	}

	log.Printf("crawling file %s (%s)", args.Hash, args.Name)

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
