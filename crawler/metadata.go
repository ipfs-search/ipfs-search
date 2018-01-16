package crawler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type metadata map[string]interface{}

// getTika requests IPFS path from IPFS-TIKA and writes returned metadata
func (c *Crawler) getTika(path string, m *metadata) error {
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
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return err
	}

	return err
}

// getMatadata sets metdata for file with args or returns error
func (c *Crawler) getMetadata(args *Args, m *metadata) error {
	var err error

	if args.Size > 0 {
		if args.Size > c.Config.MetadataMaxSize {
			// Fail hard for really large files, for now
			return fmt.Errorf("%s (%s) too large, not indexing (for now)", args.Hash, args.Name)
		}

		path := filenameURL(args)

		tryAgain := true
		for tryAgain {
			err = c.getTika(path, m)

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

	return nil
}
