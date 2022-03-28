package batchinggetter

import (
	"context"
	"io"
	"io/ioutil"
	"strings"
	// "net/http"
	// "testing"
	// "time"

	// "github.com/dankinder/httpmock"
	// opensearch "github.com/opensearch-project/opensearch-go"
	// "github.com/stretchr/testify/mock"
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
		Field2 string `json:"a2"`
	}
	reqresp1 reqresp

	req2     *GetRequest
	rChan2   chan GetResponse
	dst2     interface{}
	reqresp2 reqresp

	// Mock search
	// mockAPIHandler *httpmock.MockHandler
	// mockAPIServer  *httpmock.Server
	// responseHeader http.Header
}

func (s *BulkRequestTestSuite) SetupTest() {
	s.ctx = context.Background()

	s.req1 = &GetRequest{
		Index:      "test",
		Fields:     []string{"a1", "a2"},
		DocumentID: "5",
	}
	s.rChan1 = make(chan GetResponse, 1)
	s.reqresp1 = reqresp{s.req1, s.rChan1, s.dst1}

	s.req2 = &GetRequest{
		Index:      "test",
		Fields:     []string{"a1", "a2"},
		DocumentID: "7",
	}
	s.rChan2 = make(chan GetResponse, 1)
	s.reqresp2 = reqresp{s.req2, s.rChan2, s.dst2}

	// Setup mock search API
	// s.mockAPIHandler = &httpmock.MockHandler{}
	// s.mockAPIServer = httpmock.NewServer(s.mockAPIHandler)
	// s.responseHeader = http.Header{
	// 	"Content-Type": []string{"application/json"},
	// }
	// client, _ := opensearch.NewClient(opensearch.Config{
	// 	Addresses: []string{s.mockAPIServer.URL()},
	// })

	// s.expectHelloWorld()
}

func (s *BulkRequestTestSuite) TestGetSearchRequest() {

	br := newBulkRequest()
	br.add(s.reqresp1)

	sr := br.getSearchRequest()
	s.Equal(sr.Index, "test1")
	s.Equal(sr.SourceIncludes, []string{"a1", "a2"})

	body := new(strings.Builder)
	io.Copy(body, sr.Body)

	s.JSONEq(`
	{
		"query": {
			"id": {
				"values": ["5", "7")]
			}
		}
	}
	`, body.String())
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
