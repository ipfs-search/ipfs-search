package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/dokterbob/ipfs-search/indexer"
	"github.com/dokterbob/ipfs-search/queue"
	"gopkg.in/ipfs/go-ipfs-api.v1"
	"log"
	"net"
	"net/http"
	"net/url"
	// "path"
	"strings"
	"syscall"
	"time"
)

const (
	// Reconnect time in seconds
	RECONNECT_WAIT = 2
	TIKA_TIMEOUT   = 120
)

type CrawlerArgs struct {
	Hash       string
	Name       string
	Size       uint64
	ParentHash string
	ParentName string
}

type Crawler struct {
	sh *shell.Shell
	id *indexer.Indexer
	fq *queue.TaskQueue
	hq *queue.TaskQueue
}

func NewCrawler(sh *shell.Shell, id *indexer.Indexer, fq *queue.TaskQueue, hq *queue.TaskQueue) *Crawler {
	return &Crawler{
		sh: sh,
		id: id,
		fq: fq,
		hq: hq,
	}
}

func hashUrl(hash string) string {
	return fmt.Sprintf("/ipfs/%s", hash)
}

// Helper function for creating reference structure
/*
	'<hash>': {
		'references': {
			[{
				'parent_hash'
				'hash'
				'name'
			}, ]
		}
	}

	if (document_exists) {
		if (references_exists) {
			add_parent_hash to references
		} else {
			add references to document
		}
	} else {
		create document with references as only information
	}
*/
func construct_references(name string, parent_hash string, parent_name string) []map[string]interface{} {
	references := []map[string]interface{}{}

	if name != "" {
		references = []map[string]interface{}{
			{
				"name":        name,
				"parent_hash": parent_hash,
				"parent_name": parent_name,
			},
		}
	}

	return references
}

// Handle IPFS errors graceously, returns try again bool and original error
func (c Crawler) handleError(err error, hash string) (bool, error) {
	if _, ok := err.(*shell.Error); ok && strings.Contains(err.Error(), "proto") {
		// We're not recovering from protocol errors, so panic

		// Attempt to index panic to prevent re-indexing
		metadata := map[string]interface{}{
			"error": err.Error(),
		}

		c.id.IndexItem("invalid", hash, metadata)

		panic(err)
	}

	if uerr, ok := err.(*url.Error); ok {
		// URL errors

		log.Printf("URL error %v", uerr)

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

// Given a particular hash (file or directory), start crawling
func (c Crawler) CrawlHash(hash string, name string, parent_hash string, parent_name string) error {
	indexed, err := c.id.IsIndexed(hash)
	if err != nil {
		return err
	}

	if indexed {
		log.Printf("Already indexed '%s', skipping", hash)
		return nil
	}

	log.Printf("Crawling hash '%s' (%s)", hash, name)

	url := hashUrl(hash)

	var list *shell.UnixLsObject

	try_again := true
	for try_again {
		list, err = c.sh.FileList(url)

		try_again, err = c.handleError(err, hash)

		if try_again {
			log.Printf("Retrying in %d seconds", RECONNECT_WAIT)
			time.Sleep(RECONNECT_WAIT * time.Duration(time.Second))
		}
	}

	if err != nil {
		return err
	}

	switch list.Type {
	case "File":
		// Add to file crawl queue
		// Note: we're expecting no references here, see comment below
		args := CrawlerArgs{
			Hash: hash,
			Name: name,
			Size: list.Size,
		}

		err = c.fq.AddTask(args)
		if err != nil {
			// failed to send the task
			return err
		}
	case "Directory":
		// Index name and size for directory and directory items
		properties := map[string]interface{}{
			"links":      list.Links,
			"size":       list.Size,
			"references": construct_references(name, parent_hash, parent_name),
		}

		err := c.id.IndexItem("directory", hash, properties)
		if err != nil {
			return err
		}

		for _, link := range list.Links {
			args := CrawlerArgs{
				Hash:       link.Hash,
				Name:       link.Name,
				Size:       link.Size,
				ParentHash: hash,
				ParentName: name,
			}

			switch link.Type {
			case "File":
				// Add file to crawl queue
				err = c.fq.AddTask(args)
				if err != nil {
					// failed to send the task
					return err
				}

			case "Directory":
				// Add directory to crawl queue
				c.hq.AddTask(args)
				if err != nil {
					// failed to send the task
					return err
				}
			default:
				log.Printf("Type '%s' skipped for '%s'", list.Type, hash)
			}
		}
	default:
		log.Printf("Type '%s' skipped for '%s'", list.Type, hash)
	}

	log.Printf("Finished hash %s", hash)

	return nil
}

func getMetadata(path string, metadata *map[string]interface{}) error {
	const ipfs_tika_url = "http://localhost:8081"

	client := http.Client{
		Timeout: TIKA_TIMEOUT * time.Duration(time.Second),
	}

	resp, err := client.Get(ipfs_tika_url + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Undesired status '%s' from ipfs-tika.", resp.Status)
	}

	// Parse resulting JSON
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return err
	}

	return err
}

// Crawl a single object, known to be a file
func (c Crawler) CrawlFile(hash string, name string, parent_hash string, parent_name string, size uint64) error {
	indexed, err := c.id.IsIndexed(hash)
	if err != nil {
		return err
	}

	if indexed {
		log.Printf("Already indexed '%s', skipping", hash)
		return nil
	}

	log.Printf("Crawling file %s (%s)\n", hash, name)

	metadata := make(map[string]interface{})

	if size > 0 {
		if size > 10*1024*1024 {
			// Fail hard for really large files, for now
			return fmt.Errorf("%s (%s) too large, not indexing (for now).", hash, name)
		}

		var path string
		if name != "" && parent_hash != "" {
			path = fmt.Sprintf("/ipfs/%s/%s", parent_hash, name)
		} else {
			path = fmt.Sprintf("/ipfs/%s", hash)
		}

		try_again := true
		for try_again {
			err = getMetadata(path, &metadata)

			try_again, err = c.handleError(err, hash)

			if try_again {
				log.Printf("Retrying in %d seconds", RECONNECT_WAIT)
				time.Sleep(RECONNECT_WAIT * time.Duration(time.Second))
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
					args := CrawlerArgs{
						Hash:       link.Hash,
						Name:       link.Name,
						Size:       link.Size,
						ParentHash: hash,
						ParentName: name,
					}

				}
			}
		*/
	}

	metadata["size"] = size
	metadata["references"] = construct_references(name, parent_hash, parent_name)

	err = c.id.IndexItem("file", hash, metadata)
	if err != nil {
		return err
	}

	log.Printf("Finished file %s", hash)

	return nil
}
