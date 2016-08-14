package crawler

import (
	"fmt"
	"github.com/dokterbob/ipfs-search/indexer"
	"gopkg.in/ipfs/go-ipfs-api.v1"
)

type Crawler struct {
	sh *shell.Shell
	id *indexer.Indexer
}

func NewCrawler(sh *shell.Shell, id *indexer.Indexer) *Crawler {
	c := new(Crawler)
	c.sh = sh
	c.id = id
	return c
}

// Given a particular hash, start crawling
func (c Crawler) CrawlHash(hash string) error {
	fmt.Printf("Crawling hash %s\n", hash)

	url := fmt.Sprintf("/ipfs/%s", hash)

	list, err := c.sh.FileList(url)
	if err != nil {
		return err
	}

	switch list.Type {
	case "File":
		// Add to file crawl queue
		err := c.CrawlFile(hash)
		if err != nil {
			return err
		}
	case "Directory":
		// Index name and size for items
		c.id.IndexDirectory(list)

		for _, link := range list.Links {
			switch link.Type {
			case "File":
				// Add file to crawl queue
				err := c.CrawlFile(link.Hash)

				if err != nil {
					return err
				}
			case "Directory":
				// Add directory to crawl queue
				err := c.CrawlHash(link.Hash)

				if err != nil {
					return err
				}
			}
		}
	default:
		fmt.Printf("Type not '%s' skipped for '%s'", list.Type, hash)
	}

	return nil
}

// Crawl a single object, known to be a file
func (c Crawler) CrawlFile(hash string) error {
	fmt.Printf("Crawling file %s\n", hash)

	return nil
}
