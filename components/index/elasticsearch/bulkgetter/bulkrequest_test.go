package bulkgetter

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/stretchr/testify/suite"
)

type BulkRequestTestSuite struct {
	suite.Suite
	ctx context.Context
	br  bulkRequest

	req1   *GetRequest
	rChan1 chan GetResponse
	dst1   struct {
		Field1 string `json:"a1"`
		Field2 int    `json:"a2"`
	}
	reqresp1 reqresp

	req2     *GetRequest
	rChan2   chan GetResponse
	dst2     interface{}
	reqresp2 reqresp
}

func (s *BulkRequestTestSuite) SetupTest() {
	s.ctx = context.Background()

	s.req1 = &GetRequest{
		Index:      "test",
		Fields:     []string{"a1", "a2"},
		DocumentID: "5",
	}
	s.rChan1 = make(chan GetResponse, 1)
	s.reqresp1 = reqresp{s.req1, s.rChan1, &s.dst1}

	s.req2 = &GetRequest{
		Index:      "test",
		Fields:     []string{"a1", "a2"},
		DocumentID: "7",
	}
	s.rChan2 = make(chan GetResponse, 1)
	s.reqresp2 = reqresp{s.req2, s.rChan2, &s.dst2}
}

func (s *BulkRequestTestSuite) TestGetSearchRequest() {
	br := newBulkRequest()
	br.add(s.reqresp1)
	br.add(s.reqresp2)

	sr := br.getSearchRequest()
	s.Equal([]string{"test"}, sr.Index)
	s.Equal([]string{"a1", "a2"}, sr.SourceIncludes)

	var res struct {
		Query struct {
			Ids struct {
				Values []string
			}
		}
		Sort []string
	}

	err := json.NewDecoder(sr.Body).Decode(&res)
	s.NoError(err)

	values := res.Query.Ids.Values
	s.Contains(values, "5")
	s.Contains(values, "7")

	sort := res.Sort
	s.Equal("_doc", sort[0])
}

func (s *BulkRequestTestSuite) TestProcessResponseFound() {
	br := newBulkRequest()
	br.add(s.reqresp1)
	br.add(s.reqresp2)

	respStr := `{
	  "took": 3,
	  "timed_out": false,
	  "_shards": {
	    "total": 1,
	    "successful": 1,
	    "skipped": 0,
	    "failed": 0
	  },
	  "hits": {
	    "total": {
	      "value": 1,
	      "relation": "eq"
	    },
	    "max_score": 1.0,
	    "hits": [
	      {
	        "_index": "test",
	        "_type": "_doc",
	        "_id": "5",
	        "_score": 1.0,
	        "_source": {
	          "a1": "kaas",
	          "a2": 15
	        }
	      }
	    ]
	  }
	}`

	resp := opensearchapi.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(respStr)),
	}

	err := br.processResponse(&resp)
	s.NoError(err)

	s.NotEmpty(s.reqresp1.resp)
	s.NotEmpty(s.reqresp2.resp)

	r1 := <-s.reqresp1.resp
	r2 := <-s.reqresp2.resp
	s.True(r1.Found)
	s.False(r2.Found)

	s.Equal("kaas", s.dst1.Field1)
	s.Equal(15, s.dst1.Field2)
}

func TestBulkRequestTestSuite(t *testing.T) {
	suite.Run(t, new(BulkRequestTestSuite))
}
