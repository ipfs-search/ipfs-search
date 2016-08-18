package crawler

import (
	"errors"
	"fmt"
	machinery "github.com/RichardKnop/machinery/v1"
	signatures "github.com/RichardKnop/machinery/v1/signatures"
	"github.com/dokterbob/ipfs-search/indexer"
	"gopkg.in/ipfs/go-ipfs-api.v1"
	"net/http"
)

type Crawler struct {
	sh  *shell.Shell
	id  *indexer.Indexer
	mac *machinery.Server
}

func NewCrawler(sh *shell.Shell, id *indexer.Indexer, mac *machinery.Server) *Crawler {
	c := new(Crawler)
	c.sh = sh
	c.id = id
	c.mac = mac
	return c
}

func hashUrl(hash string) string {
	return fmt.Sprintf("/ipfs/%s", hash)
}

// Given a particular hash, start crawling
func (c Crawler) CrawlHash(hash string) error {
	indexed, err := c.id.IsIndexed(hash)
	if err != nil {
		return err
	}

	if indexed {
		fmt.Printf("Already indexed '%s', skipping\n", hash)
		return nil
	}

	fmt.Printf("Crawling hash '%s'\n", hash)

	url := hashUrl(hash)

	list, err := c.sh.FileList(url)
	if err != nil {
		return err
	}

	switch list.Type {
	case "File":
		// Add to file crawl queue
		task := signatures.TaskSignature{
			Name: "crawl_file",
			Args: []signatures.TaskArg{
				signatures.TaskArg{
					Type:  "string",
					Value: hash,
				},
				signatures.TaskArg{
					Type:  "string",
					Value: nil,
				},
			},
		}
		_, err := c.mac.SendTask(&task)
		if err != nil {
			// failed to send the task
			return err
		}
	case "Directory":
		// Index name and size for items
		properties := map[string]interface{}{
			"links": list.Links,
		}

		c.id.IndexItem("Directory", hash, properties)

		for _, link := range list.Links {
			c.id.IndexReference(link.Type, link.Hash, link.Name, hash)

			switch link.Type {
			case "File":
				// Add file to crawl queue
				task := signatures.TaskSignature{
					Name: "crawl_file",
					Args: []signatures.TaskArg{
						signatures.TaskArg{
							Type:  "string",
							Value: link.Hash,
						},
					},
				}
				_, err := c.mac.SendTask(&task)
				if err != nil {
					// failed to send the task
					return err
				}

			case "Directory":
				// Add directory to crawl queue
				task := signatures.TaskSignature{
					Name: "crawl_hash",
					Args: []signatures.TaskArg{
						signatures.TaskArg{
							Type:  "string",
							Value: link.Hash,
						},
					},
				}
				_, err := c.mac.SendTask(&task)
				if err != nil {
					// failed to send the task
					return err
				}
			default:
				fmt.Printf("Type '%s' skipped for '%s'", list.Type, hash)
			}
		}
	default:
		fmt.Printf("Type '%s' skipped for '%s'", list.Type, hash)
	}

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
	if err != nil && err.Error() != "EOF" {
		return "", err
	}

	if numread == 0 {
		return "", errors.New("0 characters read, mime type detection failed")
	}

	// Sniffing only uses at most the first 512 bytes
	return http.DetectContentType(data), nil
}

// Crawl a single object, known to be a file
func (c Crawler) CrawlFile(hash string) error {
	fmt.Printf("Crawling file %s\n", hash)

	mimetype, err := c.getMimeType(hash)
	if err != nil {
		return err
	}

	properties := map[string]interface{}{
		"mimetype": mimetype,
	}

	c.id.IndexItem("File", hash, properties)

	return nil
}
