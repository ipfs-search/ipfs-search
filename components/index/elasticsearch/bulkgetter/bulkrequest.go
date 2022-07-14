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
)

// ErrHTTP represents non-404 errors in HTTP requests.
var ErrHTTP = errors.New("HTTP Error")

type bulkRequest struct {
	ctx         context.Context
	client      *opensearch.Client
	rrs         map[string]reqresp
	decodeMutex sync.Mutex
	aliases     map[string]string
}

func newBulkRequest(ctx context.Context, client *opensearch.Client, size int) *bulkRequest {
	if ctx == nil {
		panic("required context is nil")
	}
	return &bulkRequest{
		ctx:     ctx,
		client:  client,
		rrs:     make(map[string]reqresp, size),
		aliases: make(map[string]string),
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

type aliasesResponse map[string]struct {
	Aliases map[string]struct{} `json:"aliases"`
}

func (r *bulkRequest) getAliases(indexOrAlias string) (aliasesResponse, error) {
	response := aliasesResponse{}

	falseConst := true
	req := opensearchapi.IndicesGetAliasRequest{
		Index:           []string{indexOrAlias},
		AllowNoIndices:  &falseConst,
		ExpandWildcards: "none",
	}

	res, err := req.Do(r.ctx, r.client)
	if err != nil {
		return response, fmt.Errorf("error executing request: %w", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		return response, fmt.Errorf("%w: %s", ErrHTTP, res)
	}

	err = json.NewDecoder(res.Body).Decode(&response)

	return response, err
}

func (r *bulkRequest) resolveAlias(indexOrAlias string) (string, error) {
	// GET /<index_or_alias>/_alias
	// {
	// 	"<index>": {
	// 		"aliases": {
	// 			"ipfs_directories": {}
	// 		}
	// 	}
	// }

	index, ok := r.aliases[indexOrAlias]
	if ok {
		return index, nil
	}

	response, err := r.getAliases(indexOrAlias)
	if err != nil {
		return "", err
	}

	for k := range response {
		r.aliases[indexOrAlias] = k
		return k, nil
	}

	return "", fmt.Errorf("index or alias %s not found", indexOrAlias)
}

func (r *bulkRequest) keyFromResponseDoc(doc *responseDoc) string {
	return doc.Index + doc.ID
}

func (r *bulkRequest) keyFromRR(rr reqresp) (string, error) {
	indexName, err := r.resolveAlias(rr.req.Index)
	if err != nil {
		return "", err
	}
	return indexName + rr.req.DocumentID, nil
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

	// log.Printf("Sending response to %v", rr.resp)
	// defer log.Printf("Done sending response")

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

	trueConst := true

	req := opensearchapi.MgetRequest{
		Body:       body,
		Preference: "_local",
		Realtime:   &trueConst,
	}

	return &req
}

func decodeResponse(res *opensearchapi.Response) ([]responseDoc, error) {
	// log.Printf("Decoding response to bulk GET")
	// defer log.Printf("Done decoding response to bulk GET")

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
	// log.Printf("Processing response to bulk GET")
	// defer log.Printf("Done processing response to bulk GET")

	var err error

	if res.StatusCode == 200 {
		docs, err := decodeResponse(res)
		if err != nil {
			err = fmt.Errorf("error decoding body: %w", err)
			return err
		}

		// log.Printf("Processing %d returned documents", len(docs))

		for _, d := range docs {
			key := r.keyFromResponseDoc(&d)

			// Only decode and send response when the other side is listening.
			rr, ok := r.rrs[key]
			if !ok {
				log.Printf("%+v", r.rrs)
				panic("")
				return fmt.Errorf("unknown key '%s' in response to bulk request", key)
			}
			if rr.ctx.Err() == nil {
				found, err := r.processResponseDoc(&d, key)
				r.sendResponse(key, found, err)
			} else {
				log.Printf("Not writing response from bulk get, request context cancelled.")
				close(rr.resp)
			}
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

func (r *bulkRequest) execute() error {
	log.Printf("Performing bulk GET, %d elements", len(r.rrs))

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
