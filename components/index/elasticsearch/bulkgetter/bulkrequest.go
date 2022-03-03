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

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

// ErrHTTP represents non-404 errors in HTTP requests.
var ErrHTTP = errors.New("HTTP Error")

type bulkRequest struct {
	rrs         map[string]reqresp
	decodeMutex sync.Mutex
}

func newBulkRequest(size int) *bulkRequest {
	return &bulkRequest{
		rrs: make(map[string]reqresp, size),
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

func keyFromResponseDoc(doc *responseDoc) string {
	// TODO: Resolve aliases; indexes in results are real indexes, whereas request indexes might be aliases!
	return doc.Index + doc.ID
}

func keyFromRR(rr reqresp) string {
	return rr.req.Index + rr.req.DocumentID
}

func (r *bulkRequest) add(rr reqresp) {
	r.rrs[keyFromRR(rr)] = rr
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
	// Wrap Unmarshall in mutex to prevent race conditions as dst might be shared!
	r.decodeMutex.Lock()
	defer r.decodeMutex.Unlock()

	return json.Unmarshal(src, dst)
}

// processResponseDoc returns found, error
func (r *bulkRequest) processResponseDoc(d *responseDoc, dst interface{}) (bool, error) {
	if d.Found {
		if err := r.decodeSource(d.Source, dst); err != nil {
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
			key := keyFromResponseDoc(&d)
			found, err := r.processResponseDoc(&d, r.rrs[key].dst)
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

func (r *bulkRequest) execute(ctx context.Context, client *opensearch.Client) error {
	log.Printf("Performing bulk GET, %d elements", len(r.rrs))

	res, err := r.getRequest().Do(ctx, client)
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
