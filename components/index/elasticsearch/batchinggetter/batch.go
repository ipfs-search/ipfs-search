package batchinggetter

import (
	"strings"
)

type batch map[string]map[string]bulkRequest

func getFieldsKey(fields []string) string {
	return strings.Join(fields, "")
}

func (b batch) add(rr reqresp) {
	b[getFieldsKey(rr.req.Fields)][rr.req.Index][rr.req.DocumentID] = rr
}
