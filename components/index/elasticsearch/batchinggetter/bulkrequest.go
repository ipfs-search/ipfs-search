package batchinggetter

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

type bulkRequest map[string]reqresp

func (r bulkRequest) bulkResponse(found bool, err error) {
	for _, rr := range r {
		rr.resp <- GetResponse{found, err}
		close(rr.resp)
		// Note that this does not do delete() as it should become irrelevant/unnecessary here.
	}
}

func (r bulkRequest) sendResponse(id string, found bool, err error) {
	rr := r[id]
	rr.resp <- GetResponse{found, err}
	close(rr.resp)
	delete(r, id) // Is delete the best way to do this, or setting to nil?
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

type hit struct {
	Index      string          `json:"_index"`
	DocumentID string          `json:"_id`
	Source     json.RawMessage `json:"_source"`
}

func decodeResponse(res *opensearchapi.Response) ([]hit, error) {
	response := struct {
		Hits struct {
			Hits []hit `json:"hits"`
		} `json:"hits"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Hits.Hits, nil
}

func (r bulkRequest) processResponse(res *opensearchapi.Response) error {
	switch res.StatusCode {
	case 200:
		// Found

		hits, err := decodeResponse(res)
		if err != nil {
			r.bulkResponse(false, err)
			return fmt.Errorf("error decoding body: %w", err)
		}

		for _, h := range hits {
			id := h.DocumentID

			if err := json.Unmarshal(h.Source, r[id].dst); err != nil {
				err = fmt.Errorf("error decoding source: %w", err)
				r.sendResponse(id, false, err)

				return err
			}

			// Note: this removes items from bulkRequest, so that a bulk 404 works.
			r.sendResponse(id, true, nil)
		}

	case 404:
		// None found, pass so below w can mark all remaining documents as not found.

	default:
		// return err
		panic("unexpected status from search")
	}

	r.bulkResponse(false, nil)

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

	return nil
}
