package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/dokterbob/ipfs-search/indexer"
	"github.com/dokterbob/ipfs-search/queue"
	"gopkg.in/ipfs/go-ipfs-api.v1"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const (
	// Reconnect time in seconds
	RECONNECT_WAIT = 2
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
func handleError(err error) (bool, error) {
	if _, ok := err.(*shell.Error); ok && strings.Contains(err.Error(), "proto") {
		// We're not recovering from protocol errors, so panic
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

		try_again, err = handleError(err)

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
	cmd := exec.Command("tika", "-j", path)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// Standard error to system standard error
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	// Timeout process after set time
	timer := time.AfterFunc(2*time.Minute, func() {
		log.Printf("tika timeout for '%s', killing", path)
		cmd.Process.Kill()
	})

	// Parse resulting JSON
	if err := json.NewDecoder(stdout).Decode(&metadata); err != nil {
		return err
	}

	err = cmd.Wait()
	timer.Stop()

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

	log.Printf("Crawling file %s\n", hash)

	metadata := make(map[string]interface{})

	if size > 0 {
		var path string
		if name != "" && parent_hash != "" {
			path = fmt.Sprintf("/ipfs/%s/%s", parent_hash, name)
		} else {
			path = fmt.Sprintf("/ipfs/%s", hash)
		}

		if err := getMetadata(path, &metadata); err != nil {
			return err
		}
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
