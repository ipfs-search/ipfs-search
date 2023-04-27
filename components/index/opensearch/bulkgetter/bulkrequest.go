package bulkgetter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"

	"github.com/ipfs-search/ipfs-search/components/index/opensearch/aliasresolver"
)

// ErrHTTP represents non-404 errors in HTTP requests.
var ErrHTTP = errors.New("HTTP Error")

type reqrespmap map[string]reqresp

type bulkRequest struct {
	ctx         context.Context
	client      *opensearch.Client
	rrs         reqrespmap
	decodeMutex sync.Mutex
	aResolver   aliasresolver.AliasResolver
}

func newBulkRequest(ctx context.Context, client *opensearch.Client, aliasResolver aliasresolver.AliasResolver, size int) *bulkRequest {
	if ctx == nil {
		panic("required context is nil")
	}

	return &bulkRequest{
		ctx:       ctx,
		client:    client,
		rrs:       make(reqrespmap, size),
		aResolver: aliasResolver,
	}
}

func (r *bulkRequest) sendBulkResponse(found bool, err error) {
	for _, rr := range r.rrs {
		rr.resp <- GetResponse{found, err}
		close(rr.resp)
	}
}

type responseDoc struct {
	Index  string          `json:"_index"`
	ID     string          `json:"_id"`
	Found  bool            `json:"found"`
	Source json.RawMessage `json:"_source"`
}

func (r *bulkRequest) keyFromResponseDoc(doc *responseDoc) string {
	return doc.Index + doc.ID
}

func (r *bulkRequest) keyFromRR(rr reqresp) (string, error) {
	// Resolve index back to alias
	aliasName, err := r.aResolver.GetAlias(rr.ctx, rr.req.Index)
	if err != nil {
		return "", err
	}
	return aliasName + rr.req.DocumentID, nil
}

func (r *bulkRequest) add(rr reqresp) error {
	key, err := r.keyFromRR(rr)
	if err != nil {
		return err
	}

	r.rrs[key] = rr

	return nil
}

func (r *bulkRequest) sendResponse(key string, found bool, err error) {
	rr, keyFound := r.rrs[key]

	if !keyFound {
		panic(fmt.Sprintf("Key %s not found in reqresp %v.", key, r.rrs))
	}

	if rr.resp == nil {
		panic(fmt.Sprintf("Invalid value for response channel for reqresp %v", rr))
	}

	if debug {
		log.Printf("bulkrequest: Sending response for %v", &rr)
		defer log.Printf("bulkrequest: Done sending response")
	}

	rr.resp <- GetResponse{found, err}
	close(rr.resp)
}

func (r *bulkRequest) getReqBody() io.Reader {
	type source struct {
		Include []string `json:"include"`
	}

	type doc struct {
		Index  string `json:"_index"`
		ID     string `json:"_id"`
		Source source `json:"_source"`
	}

	docs := make([]doc, len(r.rrs))

	i := 0
	for _, rr := range r.rrs {
		docs[i] = doc{
			Index: rr.req.Index,
			ID:    rr.req.DocumentID,
			Source: source{
				rr.req.Fields,
			},
		}

		i++
	}

	bodyStruct := struct {
		Docs []doc `json:"docs"`
	}{docs}

	var buffer bytes.Buffer

	e := json.NewEncoder(io.Writer(&buffer))
	if err := e.Encode(bodyStruct); err != nil {
		panic("Error generating MGET request body.")
	}

	return io.Reader(&buffer)
}

func (r *bulkRequest) getRequest() *opensearchapi.MgetRequest {
	body := r.getReqBody()

	req := opensearchapi.MgetRequest{
		Body: body,
	}

	return &req
}

func decodeResponse(res *opensearchapi.Response) ([]responseDoc, error) {
	if debug {
		log.Printf("bulkrequest: Decoding response to bulk GET")
		defer log.Printf("bulkrequest: Done decoding response to bulk GET")
	}

	response := struct {
		Docs []responseDoc `json:"docs"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Docs, nil
}

func (r *bulkRequest) decodeSource(src json.RawMessage, dst interface{}) error {
	// Wrap Unmarshall in mutex to prevent race conditions as dst may be shared!
	r.decodeMutex.Lock()
	defer r.decodeMutex.Unlock()

	return json.Unmarshal(src, dst)
}

// processResponseDoc returns found, error
func (r *bulkRequest) processResponseDoc(d *responseDoc, key string) (bool, error) {
	// Only decode and send response when the other side is listening.
	rr, ok := r.rrs[key]
	if !ok {
		// Panic, this is a proper bug.
		panic(fmt.Sprintf("unknown key '%s' in response to bulk request", key))
	}

	if err := rr.ctx.Err(); err != nil {
		if debug {
			log.Printf("bulkrequest: Not writing response from bulk get, request context cancelled.")
		}

		return false, err

	}

	// Context still open
	if d.Found {
		if err := r.decodeSource(d.Source, r.rrs[key].dst); err != nil {
			err = fmt.Errorf("error decoding source: %w", err)
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (r *bulkRequest) processResponse(res *opensearchapi.Response) error {
	if debug {
		log.Printf("bulkrequest: processing response to bulk GET")
		defer log.Printf("bulkrequest: done processing response to bulk GET")
	}

	var err error

	if res.StatusCode == 200 {
		docs, err := decodeResponse(res)
		if err != nil {
			err = fmt.Errorf("error decoding body: %w", err)
			return err
		}

		if debug {
			log.Printf("bulkrequest: Processing %d returned documents", len(docs))
		}

		for _, d := range docs {
			key := r.keyFromResponseDoc(&d)
			found, err := r.processResponseDoc(&d, key)
			r.sendResponse(key, found, err)
		}

		return nil
	}

	// Non-200 status codes signify an error
	if res.IsError() {
		err = fmt.Errorf("%w: %s", ErrHTTP, res)
	} else {
		err = fmt.Errorf("Unexpected HTTP return code: %d", res.StatusCode)
	}

	return err
}

// removeCanceled removes items from query if they're canceled before the request
func (r *bulkRequest) removeCanceled() {
	removed := 0

	for key, rr := range r.rrs {
		if err := rr.ctx.Err(); err != nil {
			if debug {
				log.Printf("bulkrequest: request canceled, removing %v", &rr)
				removed++
			}

			// Send response, cleaning up resources.
			r.sendResponse(key, false, err)

			// Delete request, preventing it from being sent.
			delete(r.rrs, key)
		}
	}

	if debug && removed > 0 {
		log.Printf("bulkrequest: removed %d canceled requests", removed)
	}
}

func (r *bulkRequest) execute() error {
	r.removeCanceled()

	if len(r.rrs) == 0 {
		if debug {
			log.Printf("bulkrequest: empty request map, not sending request")
		}

		return nil
	}

	if debug {
		log.Printf("bulkrequest: performing bulk GET, %d elements", len(r.rrs))
	}

	res, err := r.getRequest().Do(r.ctx, r.client)
	if err != nil {
		err = fmt.Errorf("error executing request: %w", err)
		r.sendBulkResponse(false, err)
		return err
	}

	defer res.Body.Close()

	if err = r.processResponse(res); err != nil {
		r.sendBulkResponse(false, err)
		return err
	}

	return nil
}
