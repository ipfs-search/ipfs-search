package batchinggetter

import (
	"context"
	"strings"

	"github.com/opensearch-project/opensearch-go"
)

type batch map[string]map[string]bulkRequest

func getFieldsKey(fields []string) string {
	return strings.Join(fields, "")
}

func (b batch) add(rr reqresp) {
	b[getFieldsKey(rr.req.Fields)][rr.req.Index][rr.req.DocumentID] = rr
}

func (b batch) execute(ctx context.Context, client *opensearch.Client) error {
	for _, indexes := range b {
		for _, r := range indexes {
			err := r.performBulkRequest(ctx, client)
			if err != nil {
				// Note: this will terminate batch on first error in request.
				return err
			}
		}
	}

	return nil
}
