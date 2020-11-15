package tika

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ipfs-search/ipfs-search/extractor"
	t "github.com/ipfs-search/ipfs-search/types"
)

// Extractor extracts metadata using the ipfs-tika server.
type Extractor struct {
	config *Config
	client http.Client
}

// retryingGet is an infinitely retrying GET on intermittent errors (e.g. server goes)
// TODO: Replace by proper circuit breakers.
func (e *Extractor) retryingGet(ctx context.Context, url string) (resp *http.Response, err error) {
	retries := 0

	for {
		log.Printf("Fetching metadata from '%s'", url)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			// Errors here are programming errors.
			panic("error creating request")
		}

		resp, err = e.client.Do(req)

		// TODO: This is probably a sensible update to go, which might simplify
		// shouldRetry - but better to have tracing infra in place before we go there.
		//
		// Any returned error will be of type *url.Error. The url.Error value's Timeout
		// method will report true if request timed out or was canceled.
		// Ref: https://golang.org/pkg/net/http/#Client.Do

		if err == nil {
			// Success, we're done here.
			return resp, nil
		}

		if !shouldRetry(err) {
			// Fatal error
			return nil, err
		}

		retries++

		log.Printf("Retrying (%d) in %s", retries, e.config.RetryWait)
		time.Sleep(e.config.RetryWait)
	}
}

func (e *Extractor) getExtractURL(r t.ReferencedResource) string {
	return e.config.TikaServerURL + r.GatewayPath()
}

// Extract metadata from a (potentially) referenced resource, updating
// Metadata or returning an error.
func (e *Extractor) Extract(ctx context.Context, r t.ReferencedResource, m t.Metadata) error {
	resp, err := e.retryingGet(ctx, e.getExtractURL(r))

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

	// TODO
	// Check for IPFS links in urls extracted from resource
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

	return nil
}

func getClient(config *Config) http.Client {
	// TODO: Get more advanced client with circuit breaking etc. over manual
	// retrying get etc.
	// Ref: https://github.com/gojek/heimdall#creating-a-hystrix-like-circuit-breaker
	return http.Client{
		Timeout: config.RequestTimeout,
	}
}

// New returns a new Tika extractor.
func New(config *Config) extractor.Extractor {
	return &Extractor{
		config: config,
		client: getClient(config),
	}
}
