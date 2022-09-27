package bulkgetter

import (
	"context"
)

type reqresp struct {
	ctx  context.Context
	req  *GetRequest
	resp chan GetResponse
	dst  interface{}
}
