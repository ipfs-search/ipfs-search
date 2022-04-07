package bulkgetter

type reqresp struct {
	req  *GetRequest
	resp chan GetResponse
	dst  interface{}
}
