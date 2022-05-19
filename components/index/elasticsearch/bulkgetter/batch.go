package bulkgetter

import (
	"context"
	"strings"

	"github.com/opensearch-project/opensearch-go"
)

// TODO: Merge batch and bulkrequest, using multi-index MGET.

// Ref: https://opensearch.org/docs/latest/opensearch/rest-api/document-apis/multi-get/
//
// GET _mget
// {
//   "docs": [
//   {
//     "_index": "sample-index1",
//     "_id": "1"
//   },
//   {
//     "_index": "sample-index2",
//     "_id": "1",
//     "_source": {
//       "include": ["Length"]
//     }
//   }
//   ]
// }
//
// https://pkg.go.dev/github.com/opensearch-project/opensearch-go@v1.1.0/opensearchapi#Mget


type batch map[string]bulkRequest

func getKey(rr reqresp) string {
	return strings.Join(rr.req.Fields, "") + rr.req.Index
}

func newBatch(size int) batch {
	return make(batch, size)
}

func (b batch) add(rr reqresp) {
	if b[getKey(rr)] == nil {
		b[getKey(rr)] = newBulkRequest()
	}

	b[getKey(rr)].add(rr)
}

func (b batch) execute(ctx context.Context, client *opensearch.Client) error {
	for _, br := range b {
		if err := br.performBulkRequest(ctx, client); err != nil {
			// Note: this will terminate batch on first error in request.
			return err
		}
	}

	return nil
}
