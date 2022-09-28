package bulkgetter

import (
	"context"
	"fmt"
)

type reqresp struct {
	ctx  context.Context
	req  *GetRequest
	resp chan GetResponse
	dst  interface{}
}

func (rr *reqresp) String() string {
	return fmt.Sprintf("reqresp: %v", rr.req)
}
