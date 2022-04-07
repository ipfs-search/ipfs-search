package bulkgetter

import (
	"context"
	"strings"

	"github.com/opensearch-project/opensearch-go"
)

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
