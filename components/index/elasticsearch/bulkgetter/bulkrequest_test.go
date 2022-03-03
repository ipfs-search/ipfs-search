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
		Index:      "test1",
		Fields:     []string{"a1", "a2"},
		DocumentID: "5",
	}
	s.rChan1 = make(chan GetResponse, 1)
	s.reqresp1 = reqresp{s.req1, s.rChan1, &s.dst1}

	s.req2 = &GetRequest{
		Index:      "test2",
		Fields:     []string{"b"},
		DocumentID: "7",
	}
	s.rChan2 = make(chan GetResponse, 1)
	s.reqresp2 = reqresp{s.req2, s.rChan2, &s.dst2}
}

func (s *BulkRequestTestSuite) TestGetRequest() {
	br := newBulkRequest(2)
	br.add(s.reqresp1)
	br.add(s.reqresp2)

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

	err := json.NewDecoder(r.Body).Decode(&bodyStruct)
	s.NoError(err)

	s.Equal(s.req1.Index, bodyStruct.Docs[0].Index)
	s.Equal(s.req1.DocumentID, bodyStruct.Docs[0].ID)
	s.Equal(s.req2.Index, bodyStruct.Docs[1].Index)
	s.Equal(s.req2.DocumentID, bodyStruct.Docs[1].ID)
}

func (s *BulkRequestTestSuite) TestProcessResponseFound() {
	br := newBulkRequest(2)
	br.add(s.reqresp1)
	br.add(s.reqresp2)

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
