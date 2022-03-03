package bulkgetter

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/dankinder/httpmock"
	opensearch "github.com/opensearch-project/opensearch-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BulkGetterSuite struct {
	suite.Suite
	ctx context.Context
	cfg Config
	bg  *BulkGetter

	// Mock search
	mockAPIHandler *httpmock.MockHandler
	mockAPIServer  *httpmock.Server
	responseHeader http.Header
}

func (s *BulkGetterSuite) expectHelloWorld() {
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

func (s *BulkGetterSuite) SetupTest() {
	s.ctx = context.Background()

	// Setup mock search API
	s.mockAPIHandler = &httpmock.MockHandler{}
	s.mockAPIServer = httpmock.NewServer(s.mockAPIHandler)
	s.responseHeader = http.Header{
		"Content-Type": []string{"application/json"},
	}
	client, _ := opensearch.NewClient(opensearch.Config{
		Addresses: []string{s.mockAPIServer.URL()},
	})

	s.expectHelloWorld()

	// Setup batching getter
	s.cfg = Config{
		Client:       client,
		BatchTimeout: time.Millisecond,
		BatchSize:    4,
	}
	s.bg = New(s.cfg)
}

func (s *BulkGetterSuite) TestGet() {
	req := GetRequest{}
	dst := struct{}{}
	resp := s.bg.Get(s.ctx, &req, &dst)

	s.Empty(resp)
}

func (s *BulkGetterSuite) TestProcessBatchContextCancel() {
	ctx, cancel := context.WithCancel(s.ctx)
	cancel()

	err := s.bg.processBatch(ctx)
	s.ErrorIs(err, context.Canceled)
}

func (s *BulkGetterSuite) TestProcessBatchTimeout() {
	t := time.Now()

	err := s.bg.processBatch(s.ctx)
	s.NoError(err)

	s.True(time.Now().After(t.Add(s.cfg.BatchTimeout)))
}

func (s *BulkGetterSuite) TestPopulateBatch() {
	var dst interface{}
	queue := make(chan reqresp, 4)
	rChan := make(chan GetResponse, 1)

	reqresp1 := reqresp{&GetRequest{
		Index:      "index1",
		Fields:     []string{"f1", "f2"},
		DocumentID: "1",
	}, rChan, dst}

	// Same everything, should be grouped
	reqresp2 := reqresp{&GetRequest{
		Index:      "index1",
		Fields:     []string{"f1", "f2"},
		DocumentID: "2",
	}, rChan, dst}

	// Different index
	reqresp3 := reqresp{&GetRequest{
		Index:      "index2",
		Fields:     []string{"f1", "f2"},
		DocumentID: "3",
	}, rChan, dst}

	// Different fields
	reqresp4 := reqresp{&GetRequest{
		Index:      "index1",
		Fields:     []string{"f1", "f3"},
		DocumentID: "4",
	}, rChan, dst}

	// Async; prevent buffer lock up.
	go func() {
		queue <- reqresp1
		queue <- reqresp2
		queue <- reqresp3
		queue <- reqresp4

		// Same everything, but outside batch range, so different request
		queue <- reqresp1
	}()

	b, err := s.bg.populateBatch(s.ctx, queue)

	s.Len(b.rrs, 4)
	s.NoError(err)

	s.Equal(reqresp1, b.rrs[keyFromRR(reqresp1)])
}

// TestProcessBatch is an integration test.
func (s *BulkGetterSuite) TestProcessBatch() {
	testFound := []byte(`{
	  "docs": [
	    {
	      "_index": "test_index_1",
	      "_id": "1",
	      "_version": 4,
	      "_seq_no": 5,
	      "_primary_term": 19,
	      "found": true,
	      "_source": {
	        "field1": "kaas",
	        "field2": 15
	      }
	    },
	    {
	      "_index": "test_index_2",
	      "_id": "2",
	      "_version": 1,
	      "_seq_no": 6,
	      "_primary_term": 19,
	      "found": false
	    }
	  ]
	}`)

	testURL := "/_mget?preference=_local&realtime=true"
	s.mockAPIHandler.
		On("Handle", "POST", testURL, mock.Anything).
		Return(httpmock.Response{
			Body: testFound,
		}).
		Once()

	type testType struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	dst1 := testType{}

	// Expected: found
	req1 := GetRequest{
		Index:      "test_index_1",
		DocumentID: "1",
		Fields:     []string{"field1", "field2"},
	}

	resp1 := s.bg.Get(s.ctx, &req1, &dst1)

	// Expected: not found
	dst2 := testType{}

	req2 := GetRequest{
		Index:      "test_index_2",
		DocumentID: "2",
		Fields:     []string{"field1", "field2"},
	}

	go func() {
		s.NoError(s.bg.Work(s.ctx))
	}()

	resp2 := s.bg.Get(s.ctx, &req2, &dst2)

	r1 := <-resp1
	s.True(r1.Found)
	s.Equal(dst1.Field1, "kaas")
	s.Equal(dst1.Field2, 15)

	r2 := <-resp2
	s.False(r2.Found)
	s.Empty(dst2)
}

func TestBulkGetterSuite(t *testing.T) {
	suite.Run(t, new(BulkGetterSuite))
}
