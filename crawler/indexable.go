package crawler

import (
	"net"
	"net/url"
	// "path"
	"fmt"
	"github.com/ipfs-search/ipfs-search/indexer"
	"github.com/ipfs/go-ipfs-api"
	"log"
	"strings"
	"syscall"
	"time"
)

// Indexable consists of args with a Crawler
type Indexable struct {
	*Crawler
	*Args
}

// skipItem determines whether a particular item should not be indexed
// This holds particularly to partial content.
func (i *Indexable) skipItem() bool {
	if i.Size == i.Config.PartialSize && i.ParentHash == "" {
		log.Printf("Skipping unreferenced partial content for file %s", i.Hash)
		return true
	}

	return false
}

// handleShellError handles IPFS shell errors; returns try again bool and original error
func (i *Indexable) handleShellError(err error) (bool, error) {
	if _, ok := err.(*shell.Error); ok && strings.Contains(err.Error(), "proto") {
		// We're not recovering from protocol errors, so panic

		// Attempt to index panic to prevent re-indexing
		m := metadata{
			"error": err.Error(),
		}

		i.Indexer.IndexItem("invalid", i.Hash, m)

		panic(err)
	}

	// Different error, attempt handling as URL error
	return i.handleURLError(err)
}

// handleURLError handles HTTP errors graceously, returns try again bool and original error
func (i *Indexable) handleURLError(err error) (bool, error) {
	if uerr, ok := err.(*url.Error); ok {
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

// hashURL returns the IPFS URL for a particular hash
func (i *Indexable) hashURL() string {
	return fmt.Sprintf("/ipfs/%s", i.Hash)
}

// getFileList return list of files and/or type of item (directory/file)
func (i *Indexable) getFileList() (list *shell.UnixLsObject, err error) {
	url := i.hashURL()

	tryAgain := true
	for tryAgain {
		list, err = i.Shell.FileList(url)

		tryAgain, err = i.handleShellError(err)

		if tryAgain {
			log.Printf("Retrying in %s", i.Config.RetryWait)
			time.Sleep(i.Config.RetryWait)
		}
	}

	return
}

// queueList queues any items in a given list/directory
func (i *Indexable) queueList(list *shell.UnixLsObject) (err error) {
	for _, link := range list.Links {
		dirArgs := &Args{
			Hash:       link.Hash,
			Name:       link.Name,
			Size:       link.Size,
			ParentHash: i.Hash,
		}

		switch link.Type {
		case "File":
			// Add file to crawl queue
			err = i.FileQueue.AddTask(dirArgs)
		case "Directory":
			// Add directory to crawl queue
			err = i.HashQueue.AddTask(dirArgs)
		default:
			log.Printf("Type '%s' skipped for '%s'", link.Type, i.Hash)
		}
	}

	return
}

// processList processes a file listing
func (i *Indexable) processList(list *shell.UnixLsObject, references []indexer.Reference) (err error) {
	switch list.Type {
	case "File":
		// Add to file crawl queue
		fileArgs := Args{
			Hash:       i.Hash,
			Name:       i.Name,
			Size:       list.Size,
			ParentHash: i.ParentHash,
		}

		err = i.FileQueue.AddTask(fileArgs)
	case "Directory":
		// Queue indexing of linked items
		err = i.queueList(list)
		if err != nil {
			return err
		}

		// Index name and size for directory and directory items
		properties := metadata{
			"links":      list.Links,
			"size":       list.Size,
			"references": references,
		}

		err = i.Indexer.IndexItem("directory", i.Hash, properties)
	default:
		log.Printf("Type '%s' skipped for '%s'", list.Type, i.Hash)
	}

	return
}

// updateItem updates references (and later also last seen date)
func (i *Indexable) updateItem(references []indexer.Reference) error {
	properties := metadata{
		"references": references,
	}

	return i.Indexer.IndexItem(itemType, i.Hash, properties)
}

// CrawlHash crawls a particular hash (file or directory)
func (i *Indexable) CrawlHash() error {
	if i.skipItem() {
		return nil
	}

	references, alreadyIndexed, err := i.getReferences()
	if err != nil {
		return err
	}

	if alreadyIndexed {
		return i.updateItem(references)
	}

	log.Printf("crawling hash '%s' (%s)", i.Hash, i.Name)

	list, err := i.getFileList()
	if err != nil {
		return err
	}

	err = i.processList(list, references)
	if err != nil {
		return err
	}

	log.Printf("Finished hash %s", i.Hash)

	return nil
}

// CrawlFile crawls a single object, known to be a file
func (i *Indexable) CrawlFile() error {
	if i.skipItem() {
		return nil
	}

	references, alreadyIndexed, err := i.getReferences()

	if err != nil {
		return err
	}

	if alreadyIndexed {
		return i.updateItem(references)
	}

	log.Printf("crawling file %s (%s)", i.Hash, i.Name)

	m := make(metadata)
	i.getMetadata(&m)
	if err != nil {
		return err
	}

	// Add previously found references now
	m["size"] = i.Size
	m["references"] = references

	err = i.Indexer.IndexItem("file", i.Hash, m)
	if err != nil {
		return err
	}

	log.Printf("Finished file %s", i.Hash)

	return nil
}
