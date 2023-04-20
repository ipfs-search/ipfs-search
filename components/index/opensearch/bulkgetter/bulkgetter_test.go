package bulkgetter

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dankinder/httpmock"
	"github.com/ipfs-search/ipfs-search/components/index/opensearch/testsuite"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BulkGetterSuite struct {
	testsuite.Suite

	ctx context.Context
	cfg Config
	bg  *BulkGetter
}

func (s *BulkGetterSuite) expectResolveAlias(index string) {
	testJSON := []byte(`{
		"` + index + `": {
			"aliases": {
				"` + index + `": {}
			}
		}
	}`)

	url := fmt.Sprintf("/%s/_alias?allow_no_indices=true&expand_wildcards=none", index)

	s.MockAPIHandler.
		On("Handle", "GET", url, mock.Anything).
		Return(httpmock.Response{
			Body: testJSON,
		}).
		Maybe()
}

func (s *BulkGetterSuite) SetupTest() {
	s.SetupSearchMock()

	s.ctx = context.Background()

	client, _ := opensearch.NewClient(opensearch.Config{
		Addresses: []string{s.MockAPIServer.URL()},
	})

	// Setup batching getter
	s.cfg = Config{
		Client:       client,
		BatchTimeout: time.Millisecond,
		BatchSize:    4,
	}
	s.bg = New(s.cfg)
}

func (s *BulkGetterSuite) TeardownTest() {
	s.TeardownSearchMock()
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

	reqresp1 := reqresp{s.ctx, &GetRequest{
		Index:      "index1",
		Fields:     []string{"f1", "f2"},
		DocumentID: "1",
	}, rChan, dst}

	// Same everything, should be grouped
	reqresp2 := reqresp{s.ctx, &GetRequest{
		Index:      "index1",
		Fields:     []string{"f1", "f2"},
		DocumentID: "2",
	}, rChan, dst}

	// Different index
	reqresp3 := reqresp{s.ctx, &GetRequest{
		Index:      "index2",
		Fields:     []string{"f1", "f2"},
		DocumentID: "3",
	}, rChan, dst}

	// Different fields
	reqresp4 := reqresp{s.ctx, &GetRequest{
		Index:      "index1",
		Fields:     []string{"f1", "f3"},
		DocumentID: "4",
	}, rChan, dst}

	s.expectResolveAlias("index1")
	s.expectResolveAlias("index2")

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

	key, _ := b.keyFromRR(reqresp1)
	s.Equal(reqresp1, b.rrs[key])
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

	testURL := "/_mget"
	s.MockAPIHandler.
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

	s.expectResolveAlias("test_index_1")
	s.expectResolveAlias("test_index_2")

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
