package tika

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/ipfs-search/ipfs-search/components/extractor"
	"github.com/ipfs-search/ipfs-search/components/protocol"
	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
	"github.com/ipfs-search/ipfs-search/utils"
)

// Extractor extracts metadata using the ipfs-tika server.
type Extractor struct {
	config   *Config
	getter   utils.HTTPBodyGetter
	protocol protocol.Protocol

	*instr.Instrumentation
}

func (e *Extractor) getExtractURL(r *t.AnnotatedResource) string {
	gwURL := e.protocol.GatewayURL(r)
	return fmt.Sprintf("%s/extract?url=%s", e.config.TikaExtractorURL, url.QueryEscape(gwURL))
}

// Extract metadata from a (potentially) referenced resource, updating
// Metadata or returning an error.
func (e *Extractor) Extract(ctx context.Context, r *t.AnnotatedResource, m interface{}) error {
	ctx, span := e.Tracer.Start(ctx, "extractor.tika.Extract")
	defer span.End()

	if err := extractor.ValidateMaxSize(ctx, r, e.config.MaxFileSize); err != nil {
		return err
	}

	// Timeout if extraction hasn't fully completed within this time.
	ctx, cancel := context.WithTimeout(ctx, e.config.RequestTimeout)
	defer cancel()

	body, err := e.getter.GetBody(ctx, e.getExtractURL(r), 200)
	if err != nil {
		return err
	}
	defer body.Close()

	// Parse resulting JSON
	if err := json.NewDecoder(body).Decode(m); err != nil {
		err := fmt.Errorf("%w: %v", t.ErrUnexpectedResponse, err)
		span.RecordError(err)
		return err
	}

	log.Printf("Got tika metadata metadata for '%v'", r)

	return nil
}

// New returns a new Tika extractor.
func New(config *Config, getter utils.HTTPBodyGetter, protocol protocol.Protocol, instr *instr.Instrumentation) extractor.Extractor {
	return &Extractor{
		config,
		getter,
		protocol,
		instr,
	}
}

// Compile-time assurance that implementation satisfies interface.
var _ extractor.Extractor = &Extractor{}
