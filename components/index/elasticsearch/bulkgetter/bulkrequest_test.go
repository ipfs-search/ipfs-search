package bulkgetter

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/dankinder/httpmock"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BulkRequestTestSuite struct {
	suite.Suite

	ctx    context.Context
	client *opensearch.Client

	// Mock search
	mockAPIHandler *httpmock.MockHandler
	mockAPIServer  *httpmock.Server
	responseHeader http.Header

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

func (s *BulkRequestTestSuite) expectHelloWorld() {
	testJSON := []byte(`{
	  "name" : "0fc08b13cdab",
	  "cluster_name" : "docker-cluster",
	  "cluster_uuid" : "T9t1q7kFRSyL15qVkIlWZQ",
	  "version" : {
	    "number" : "7.8.1",
	    "build_flavor" : "oss",
	    "build_type" : "docker",
	    "build_hash" : "b5ca9c58fb664ca8bf9e4057fc229b3396bf3a89",
	    "build_date" : "2020-07-21T16:40:44.668009Z",
	    "build_snapshot" : false,
	    "lucene_version" : "8.5.1",
	    "minimum_wire_compatibility_version" : "6.8.0",
	    "minimum_index_compatibility_version" : "6.0.0-beta1"
	  },
	  "tagline" : "You Know, for Search"
	}`)
	s.mockAPIHandler.
		On("Handle", "GET", "/", mock.Anything).
		Return(httpmock.Response{
			Body: testJSON,
		}).
		Once()

}

func (s *BulkRequestTestSuite) expectResolveAlias(from string, to string) {
	testJSON := []byte(`{
		"` + to + `": {
			"aliases": {
				"` + from + `": {}
			}
		}
	}`)

	url := fmt.Sprintf("/%s/_alias?allow_no_indices=true&expand_wildcards=none", from)

	s.mockAPIHandler.
		On("Handle", "GET", url, mock.Anything).
		Return(httpmock.Response{
			Body: testJSON,
		})
}

func (s *BulkRequestTestSuite) SetupTest() {
	s.ctx = context.Background()

	// Setup mock search API
	s.mockAPIHandler = &httpmock.MockHandler{}
	s.mockAPIServer = httpmock.NewServer(s.mockAPIHandler)
	s.responseHeader = http.Header{
		"Content-Type": []string{"application/json"},
	}
	s.client, _ = opensearch.NewClient(opensearch.Config{
		Addresses: []string{s.mockAPIServer.URL()},
	})

	s.expectHelloWorld()

	s.req1 = &GetRequest{
		Index:      "test1",
		Fields:     []string{"a1", "a2"},
		DocumentID: "5",
	}
	s.rChan1 = make(chan GetResponse, 1)
	s.reqresp1 = reqresp{s.ctx, s.req1, s.rChan1, &s.dst1}

	s.req2 = &GetRequest{
		Index:      "test2",
		Fields:     []string{"b"},
		DocumentID: "7",
	}
	s.rChan2 = make(chan GetResponse, 1)
	s.reqresp2 = reqresp{s.ctx, s.req2, s.rChan2, &s.dst2}
}

func (s *BulkRequestTestSuite) TestGetRequest() {
	s.expectResolveAlias("test1", "test1")
	s.expectResolveAlias("test2", "test2")

	br := newBulkRequest(s.ctx, s.client, 2)

	err := br.add(s.reqresp1)
	s.NoError(err)
	err = br.add(s.reqresp2)
	s.NoError(err)

	r := br.getRequest()

	s.Equal("_local", r.Preference)
	s.True(*(r.Realtime))

	type source struct {
		Include []string `json:"include"`
	}

	type doc struct {
		Index  string `json:"_index"`
		ID     string `json:"_id"`
		Source source `json:"_source"`
	}

	bodyStruct := struct {
		Docs []doc `json:"docs"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&bodyStruct)
	s.NoError(err)

	s.Equal(s.req1.Index, bodyStruct.Docs[0].Index)
	s.Equal(s.req1.DocumentID, bodyStruct.Docs[0].ID)
	s.Equal(s.req2.Index, bodyStruct.Docs[1].Index)
	s.Equal(s.req2.DocumentID, bodyStruct.Docs[1].ID)
}

func (s *BulkRequestTestSuite) TestProcessResponseFound() {
	s.expectResolveAlias("test1", "test1")
	s.expectResolveAlias("test2", "test2")

	br := newBulkRequest(s.ctx, s.client, 2)

	err := br.add(s.reqresp1)
	s.NoError(err)
	err = br.add(s.reqresp2)
	s.NoError(err)

	respStr := `{
	  "docs": [
	    {
	      "_index": "test1",
	      "_id": "5",
	      "_version": 4,
	      "_seq_no": 5,
	      "_primary_term": 19,
	      "found": true,
	      "_source": {
	        "a1": "kaas",
	        "a2": 15
	      }
	    },
	    {
	      "_index": "test2",
	      "_id": "7",
	      "_version": 1,
	      "_seq_no": 6,
	      "_primary_term": 19,
	      "found": false
	    }
	  ]
	}`

	resp := opensearchapi.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(respStr)),
	}

	err = br.processResponse(&resp)
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

func (s *BulkRequestTestSuite) TestResolveIndex() {
	s.expectResolveAlias("test1", "actual_index")

	br := newBulkRequest(s.ctx, s.client, 1)
	s.NoError(br.add(s.reqresp1))

	respStr := `{
	  "docs": [
	    {
	      "_index": "actual_index",
	      "_id": "5",
	      "_version": 4,
	      "_seq_no": 5,
	      "_primary_term": 19,
	      "found": true,
	      "_source": {
	        "a1": "kaas",
	        "a2": 15
	      }
	    }
	  ]
	}`

	resp := opensearchapi.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(respStr)),
	}

	s.NoError(br.processResponse(&resp))

	s.NotEmpty(s.reqresp1.resp)
}

func TestBulkRequestTestSuite(t *testing.T) {
	suite.Run(t, new(BulkRequestTestSuite))
}
