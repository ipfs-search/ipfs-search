package tika

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"

	"github.com/ipfs-search/ipfs-search/extractor"
	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs-search/ipfs-search/protocol"
	t "github.com/ipfs-search/ipfs-search/types"
)

// Extractor extracts metadata using the ipfs-tika server.
type Extractor struct {
	config   *Config
	client   *http.Client
	protocol protocol.Protocol

	*instr.Instrumentation
}

func (e *Extractor) get(ctx context.Context, url string) (resp *http.Response, err error) {
	// Temporarily disabled due to bug - the connection needs to be open until the response body has been read!
	// ctx, cancel := context.WithTimeout(ctx, e.config.RequestTimeout)
	// defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		// Errors here are programming errors.
		panic(fmt.Sprintf("creating request: %s", err))
	}

	return e.client.Do(req)
}

func (e *Extractor) getExtractURL(r *t.AnnotatedResource) string {
	// TODO: This should be TIKAURL?url=GATEWAYURL (or something similar)
	gwURL := e.protocol.GatewayURL(r)
	u, err := url.Parse(gwURL)
	if err != nil {
		panic(fmt.Sprintf("unexpected parsing error generating URL: %v", err))
	}
	return e.config.TikaServerURL + u.EscapedPath()
}

// Extract metadata from a (potentially) referenced resource, updating
// Metadata or returning an error.
func (e *Extractor) Extract(ctx context.Context, r *t.AnnotatedResource, m interface{}) error {
	ctx, span := e.Tracer.Start(ctx, "extractor.tika.Extract",
		trace.WithAttributes(label.String("cid", r.ID)),
	)
	defer span.End()

	resp, err := e.get(ctx, e.getExtractURL(r))
	if err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("unexpected status '%s' from ipfs-tika", resp.Status)
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		return err
	}

	// Parse resulting JSON
	if err := json.NewDecoder(resp.Body).Decode(m); err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		return err
	}

	log.Printf("Got metadata metadata for '%v'", r)

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

// New returns a new Tika extractor.
func New(config *Config, client *http.Client, protocol protocol.Protocol, instr *instr.Instrumentation) extractor.Extractor {
	return &Extractor{
		config,
		client,
		protocol,
		instr,
	}
}

// Compile-time assurance that implementation satisfies interface.
var _ extractor.Extractor = &Extractor{}
