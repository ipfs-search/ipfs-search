package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

// func (i *Index) Get(ctx context.Context, id string, dst interface{}, fields ...string) (bool, error) {

// GetRequest represents an item to GET.
type GetRequest struct {
	Index      string
	DocumentID string
	Fields     []string
}

// GetResponse represents the response from a GetRequest.
type GetResponse struct {
	Found bool
	Error error
}

// AsyncGetter is an interface to allow for asynchronous getting.
type AsyncGetter interface {
	Get(context.Context, GetRequest) <-chan GetResponse
}

// BatchingGetterConfig provides configuration for a BatchingGetter.
type BatchingGetterConfig struct {
	Client    *opensearch.Client
	BatchSize int
}

type reqresp struct {
	req  *GetRequest
	resp chan GetResponse
	dst  interface{}
}

type bulkRequest map[string]reqresp

func (r bulkRequest) bulkResponse(found bool, err error) {
	for _, rr := range r {
		rr.resp <- GetResponse{found, err}
	}
}

func (r bulkRequest) getSearchRequest() *opensearchapi.SearchRequest {
	// Populate cids and get original fields value
	var (
		fields, index []string
		i             int
	)

	ids := make([]string, len(r))

	for id, rr := range r {
		ids[i] = id

		if i == 0 {
			fields = rr.req.Fields
			index = []string{rr.req.Index}
		}

		i++
	}

	body := getReqBody(ids)

	req := opensearchapi.SearchRequest{
		Index:          index,
		SourceIncludes: fields,
		Body:           strings.NewReader(body),
		// Preference:     "_local",
	}

	return &req
}

func getReqBody(ids []string) string {
	return `
	{
		"query": {
			"id": {
				"values": [` + strings.Join(ids, ", ") + `]
			}
		}
	}
	`
}

func (r bulkRequest) processResponse(res *opensearchapi.Response) error {
	switch res.StatusCode {
	case 200:
		// Found

		type hit struct {
			Index      string          `json:"_index"`
			DocumentID string          `json:"_id`
			Source     json.RawMessage `json:"_source"`
		}

		response := struct {
			Hits struct {
				Hits []hit `json:"hits"`
			} `json:"hits"`
		}{}

		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			err = fmt.Errorf("error decoding body: %w", err)
			// span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
			r.bulkResponse(false, err)
			return err
		}

		for _, hit := range response.Hits.Hits {
			// Write destination data and response for hits
			rr := r[hit.DocumentID]

			// Decode source into destination
			err := json.Unmarshal(hit.Source, rr.dst)
			if err != nil {
				err = fmt.Errorf("error decoding source: %w", err)
				// span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
				r.bulkResponse(false, err)
				return err
			}
			rr.resp <- GetResponse{true, nil}

			// Remove from map to separate found from not found
			delete(r, hit.DocumentID)
		}

	case 404:
		// None found, effectively mark all documents as not found.

	default:
		// return err
		panic("unexpected status from search")
	}

	return nil
}

func (r bulkRequest) performBulkRequest(ctx context.Context, client *opensearch.Client) error {
	// Perform search request
	res, err := r.getSearchRequest().Do(ctx, client)
	if err != nil {
		r.bulkResponse(false, err)
		return err
	}

	defer res.Body.Close()

	if err = r.processResponse(res); err != nil {
		return err
	}

	r.bulkResponse(false, nil)

	return nil
}

type batch map[string]map[string]bulkRequest

// BatchingGetter allows batching/bulk gets.
type BatchingGetter struct {
	config BatchingGetterConfig
	queue  chan reqresp
}

// NewBatchingGetter returns a new BatchingGetter, setting sensible defaults for the configuration.
func NewBatchingGetter(cfg BatchingGetterConfig) BatchingGetter {
	if cfg.Client == nil {
		cfg.Client, _ = opensearch.NewDefaultClient()
	}

	if cfg.BatchSize == 0 {
		cfg.BatchSize = 100
	}

	bg := BatchingGetter{
		config: cfg,
		queue:  make(chan reqresp, cfg.BatchSize),
	}

	return bg
}

// Get queues a single Get() for a batching get.
func (bg *BatchingGetter) Get(ctx context.Context, req *GetRequest, dst interface{}) <-chan GetResponse {
	resp := make(chan GetResponse, 1)

	bg.queue <- reqresp{req, resp, dst}

	return resp
}

func getFieldsKey(fields []string) string {
	return strings.Join(fields, "")
}

func (bg *BatchingGetter) populateBatch(ctx context.Context, queue <-chan reqresp) (batch, error) {
	var b batch

	for i := 0; i < bg.config.BatchSize; i++ {
		log.Println(i)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case rr := <-queue:
			// Add batch.Add(fields, index, documentid)
			b[getFieldsKey(rr.req.Fields)][rr.req.Index][rr.req.DocumentID] = rr
		}
	}

	return b, nil
}

func (bg *BatchingGetter) performBatch(ctx context.Context, b batch) error {
	for _, indexes := range b {
		for _, r := range indexes {
			err := r.performBulkRequest(ctx, bg.config.Client)
			if err != nil {
				// Note: this will terminate batch on errors in any requests.
				return err
			}
		}
	}

	return nil
}

// Work starts a worker to process bulk gets.
// Note: the worker will terminate on any errors.
func (bg *BatchingGetter) Work(ctx context.Context) error {
	b, err := bg.populateBatch(ctx, bg.queue)

	if err != nil {
		return err
	}

	return bg.performBatch(ctx, b)
}
