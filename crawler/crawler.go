package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/dokterbob/ipfs-search/indexer"
	"github.com/dokterbob/ipfs-search/queue"
	"gopkg.in/ipfs/go-ipfs-api.v1"
	"log"
	"os/exec"
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
	c := new(Crawler)
	c.sh = sh
	c.id = id
	c.fq = fq
	c.hq = hq
	return c
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

	list, err := c.sh.FileList(url)
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

		c.id.IndexItem("directory", hash, properties)

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
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := json.NewDecoder(stdout).Decode(&metadata); err != nil {
		return err
	}

	return cmd.Wait()
}

// Crawl a single object, known to be a file
func (c Crawler) CrawlFile(hash string, name string, parent_hash string, parent_name string, size uint64) error {
	log.Printf("Crawling file %s\n", hash)

	metadata := make(map[string]interface{})

	if size > 0 {
		var path string
		if name != "" {
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

	c.id.IndexItem("file", hash, metadata)

	log.Printf("Finished file %s", hash)

	return nil
}
