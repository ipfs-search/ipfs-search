package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/ipfs-search/ipfs-search/indexer"
	"github.com/ipfs-search/ipfs-search/queue"
	"github.com/ipfs/go-ipfs-api"
	"log"
	"net"
	"net/http"
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
		metadata := map[string]interface{}{
			"error": err.Error(),
		}

		c.Indexer.IndexItem("invalid", hash, metadata)

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

func (c *Crawler) indexReferences(hash string, name string, parentHash string) ([]indexer.Reference, bool, error) {
	var alreadyIndexed bool

	references, itemType, err := c.Indexer.GetReferences(hash)
	if err != nil {
		return nil, false, err
	}

	// TODO: Handle this more explicitly, use and detect NotFound
	if references == nil {
		alreadyIndexed = false
	} else {
		alreadyIndexed = true
	}

	references, referencesUpdated := updateReferences(references, name, parentHash)

	if alreadyIndexed {
		if referencesUpdated {
			log.Printf("Found %s, reference added: '%s' from %s", hash, name, parentHash)

			properties := map[string]interface{}{
				"references": references,
			}

			err := c.Indexer.IndexItem(itemType, hash, properties)
			if err != nil {
				return nil, false, err
			}
		} else {
			log.Printf("Found %s, references not updated.", hash)
		}
	} else if referencesUpdated {
		log.Printf("Adding %s, reference '%s' from %s", hash, name, parentHash)
	}

	return references, alreadyIndexed, nil
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

func (c *Crawler) getMetadata(path string, metadata *map[string]interface{}) error {
	client := http.Client{
		Timeout: c.Config.IpfsTikaTimeout,
	}

	resp, err := client.Get(c.Config.IpfsTikaURL + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("undesired status '%s' from ipfs-tika", resp.Status)
	}

	// Parse resulting JSON
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return err
	}

	return err
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

	metadata := make(map[string]interface{})

	if args.Size > 0 {
		if args.Size > c.Config.MetadataMaxSize {
			// Fail hard for really large files, for now
			return fmt.Errorf("%s (%s) too large, not indexing (for now)", args.Hash, args.Name)
		}

		var path string
		if args.Name != "" && args.ParentHash != "" {
			path = fmt.Sprintf("/ipfs/%s/%s", args.ParentHash, args.Name)
		} else {
			path = fmt.Sprintf("/ipfs/%s", args.Hash)
		}

		tryAgain := true
		for tryAgain {
			err = c.getMetadata(path, &metadata)

			tryAgain, err = c.handleError(err, args.Hash)

			if tryAgain {
				log.Printf("Retrying in %s", c.Config.RetryWait)
				time.Sleep(c.Config.RetryWait)
			}
		}

		if err != nil {
			return err
		}

		// Check for IPFS links in content
		/*
			for raw_url := range metadata.urls {
				url, err := URL.Parse(raw_url)

				if err != nil {
					return err
				}

				if strings.HasPrefix(url.Path, "/ipfs/") {
					// Found IPFS link!
					args := crawlerArgs{
						Hash:       link.Hash,
						Name:       link.Name,
						Size:       link.Size,
						ParentHash: hash,
					}

				}
			}
		*/
	}

	metadata["size"] = args.Size
	metadata["references"] = references

	err = c.Indexer.IndexItem("file", args.Hash, metadata)
	if err != nil {
		return err
	}

	log.Printf("Finished file %s", args.Hash)

	return nil
}
