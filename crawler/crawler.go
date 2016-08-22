package crawler

import (
	"errors"
	"fmt"
	"github.com/dokterbob/ipfs-search/indexer"
	"github.com/dokterbob/ipfs-search/queue"
	"gopkg.in/ipfs/go-ipfs-api.v1"
	"io"
	"log"
	"net/http"
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
			'<parent_hash>': {
				'name': '<name>'
			}
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
func construct_references(name string, parent_hash string, parent_name string) map[string]interface{} {
	references := map[string]interface{}{}

	if name != "" {
		references = map[string]interface{}{
			parent_hash: map[string]interface{}{
				"name":        name,
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

		c.id.IndexItem("Directory", hash, properties)

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

func (c Crawler) getMimeType(hash string) (string, error) {
	url := hashUrl(hash)
	response, err := c.sh.Cat(url)
	if err != nil {
		return "", err
	}

	defer response.Close()

	var data []byte
	data = make([]byte, 512)
	numread, err := response.Read(data)
	if err == io.EOF {
		return "", err
	}

	if numread == 0 {
		return "", errors.New("0 characters read, mime type detection failed")
	}

	// Sniffing only uses at most the first 512 bytes
	return http.DetectContentType(data), nil
}

// Crawl a single object, known to be a file
func (c Crawler) CrawlFile(hash string, name string, parent_hash string, parent_name string, size uint64) error {
	log.Printf("Crawling file %s\n", hash)

	var (
		mimetype string
		err      error
	)

	if size > 0 {
		mimetype, err = c.getMimeType(hash)
		if err != nil {
			return err
		}
	}

	properties := map[string]interface{}{
		"mimetype":   mimetype,
		"size":       size,
		"references": construct_references(name, parent_hash, parent_name),
	}

	c.id.IndexItem("File", hash, properties)

	log.Printf("Finished file %s", hash)

	return nil
}
