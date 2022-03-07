package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// Work starts a worker to process bulk gets.
func (bg *BatchingGetter) Work(ctx context.Context) error {
	var batch map[string]map[string]map[string]reqresp

	// Populate batch
	for i := 0; i < bg.config.BatchSize; i++ {
		log.Println(i)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case rr := <-bg.queue:
			// Add batch.Add(fields, index, documentid)
			batch[getFieldsKey(rr.req.Fields)][rr.req.Index][rr.req.DocumentID] = rr
		}
	}

	// Populate requests
	for fields, indexes := range batch {
		for index, CIDMap := range indexes {
			// Populate cids and get original fields value
			var reqFields []string

			cids := make([]string, len(CIDMap))
			i := 0
			for cid, rr := range CIDMap {
				cids[i] = cid

				if i == 0 {
					reqFields = rr.req.Fields
				}

				i++
			}

			// Perform search request
			req := getSearchRequest(index, reqFields, cids)
			res, err := req.Do(ctx, bg.config.Client)
			if err != nil {
				// Propagate error responses
				for _, rr := range batch[fields][index] {
					rr.resp <- GetResponse{false, err}
				}
				continue
				// return err
			}
			defer res.Body.Close()

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
					// Propagate error responses
					for _, rr := range batch[fields][index] {
						rr.resp <- GetResponse{false, err}
					}
					continue
					// return err
				}

				for _, hit := range response.Hits.Hits {
					// Write destination data and response for hits
					rr := batch[fields][index][hit.DocumentID]

					// Decode source into destination
					err = json.Unmarshal(hit.Source, rr.dst)
					if err != nil {
						err = fmt.Errorf("error decoding source: %w", err)
						// span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
						rr.resp <- GetResponse{false, err}
					} else {
						rr.resp <- GetResponse{true, nil}
					}

					// Remove from map to separate found from not found
					delete(batch[fields][index], hit.DocumentID)
				}

			case 404:
				// None found, effectively mark all documents as not found.

			default:
				// return err
				panic("unexpected status from search")
			}

			// Mark remaining documents as not found
			for _, rr := range batch[fields][index] {
				rr.resp <- GetResponse{false, nil}
			}

		}
	}

	return nil
}

func getSearchRequest(index string, fields []string, cids []string) *opensearchapi.SearchRequest {
	req := opensearchapi.SearchRequest{
		Index:          []string{index},
		SourceIncludes: fields,
		Body:           getSearchRequestBody(cids),
		// Preference:     "_local",
	}

	return &req
}

func getSearchRequestBody(ids []string) io.Reader {
	q := `
	{
		"query": {
			"id": {
				"values": [` + strings.Join(ids, ", ") + `]
			}
		}
	}
	`

	r := strings.NewReader(q)

	return r
}
