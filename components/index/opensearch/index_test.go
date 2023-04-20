package opensearch

// TODO: Test whether indexed items with omitempty are actually left out - otherwise
// non-updating references will overwrite the existing!
import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/dankinder/httpmock"
	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/ipfs-search/ipfs-search/components/index/opensearch/bulkgetter"
	"github.com/ipfs-search/ipfs-search/components/index/opensearch/testsuite"
)

type IndexTestSuite struct {
	testsuite.Suite

	ctx             context.Context
	ctxCancel       func()
	instr           *instr.Instrumentation
	mockClient      *Client
	mockAsyncGetter *bulkgetter.Mock
	responseHeader  http.Header
}

func (s *IndexTestSuite) SetupTest() {
	s.SetupSearchMock()

	s.instr = instr.New()
	s.ctx, s.ctxCancel = context.WithCancel(context.Background())

	s.mockAsyncGetter = &bulkgetter.Mock{}

	config := &ClientConfig{
		URL:   s.MockAPIServer.URL(),
		Debug: true,
	}
	s.mockClient, _ = NewClient(config, s.instr)
	s.mockClient.bulkGetter = s.mockAsyncGetter

	// Start worker
	s.mockAsyncGetter.On("Work", mock.Anything).WaitUntil(time.After(time.Second)).Return(nil).Maybe()
	go s.mockClient.Work(s.ctx)
}

func (s *IndexTestSuite) TeardownTest() {
	s.TeardownSearchMock()
}

func (s *IndexTestSuite) TestNewClient() {
	config := &ClientConfig{}
	client, err := NewClient(config, s.instr)
	s.NoError(err)
	s.NotNil(client)
}

func (s *IndexTestSuite) TestNew() {
	client, _ := NewClient(&ClientConfig{}, s.instr)
	idx := New(client, &Config{Name: "test"})
	s.NotNil(idx)
}

func (s *IndexTestSuite) TestString() {
	client, _ := NewClient(&ClientConfig{}, s.instr)
	idx := New(client, &Config{Name: "test"})
	s.Equal(fmt.Sprintf("%s", idx), "test")
}

func (s *IndexTestSuite) TestIndex() {
	idx := New(s.mockClient, &Config{Name: "test"})

	// Note whitespace here! This is NDJSON
	request := []byte(`{"create":{"_index":"test","_id":"objId"}}
{"field1":"hoi","field2":4}
`)

	response := []byte(`{
	   "took": 30,
	   "errors": false,
	   "items": [
	      {
	         "create": {
	            "_index": "test",
	            "_type": "_doc",
	            "_id": "objId",
	            "_version": 1,
	            "result": "created",
	            "_shards": {
	               "total": 2,
	               "successful": 1,
	               "failed": 0
	            },
	            "status": 201,
	            "_seq_no" : 1,
	            "_primary_term" : 2
	         }
	      }
	   ]
	}`)

	type testType struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	dst := testType{
		Field1: "hoi",
		Field2: 4,
	}

	testURL := "/_bulk"
	s.MockAPIHandler.
		On("Handle", "POST", testURL, request).
		Return(httpmock.Response{
			Body: response,
		}).
		Once()

	err := idx.Index(s.ctx, "objId", &dst)
	s.NoError(err)

	// Ensure flushing
	s.ctxCancel()
	time.Sleep(100 * time.Millisecond)

	s.MockAPIHandler.AssertExpectations(s.T())
}

func (s *IndexTestSuite) TestUpdate() {
	idx := New(s.mockClient, &Config{Name: "test"})

	// Note whitespace here! This is NDJSON
	request := []byte(`{"update":{"_index":"test","_id":"objId"}}
{"doc":{"field1":"hoi","field2":4}}
`)
	response := []byte(`{
	   "took": 30,
	   "errors": false,
	   "items": [
	      {
	         "update": {
	            "_index": "test",
	            "_type": "_doc",
	            "_id": "objId",
	            "_version": 1,
	            "result": "updated",
	            "_shards": {
	               "total": 2,
	               "successful": 1,
	               "failed": 0
	            },
	            "status": 200,
	            "_seq_no" : 1,
	            "_primary_term" : 2
	         }
	      }
	   ]
	}`)

	type testType struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	dst := testType{
		Field1: "hoi",
		Field2: 4,
	}

	testURL := "/_bulk"
	s.MockAPIHandler.
		On("Handle", "POST", testURL, request).
		Return(httpmock.Response{
			Body: response,
		}).
		Once()

	err := idx.Update(s.ctx, "objId", &dst)
	s.NoError(err)

	// Ensure flushing
	s.ctxCancel()
	time.Sleep(100 * time.Millisecond)

	s.MockAPIHandler.AssertExpectations(s.T())
}

func (s *IndexTestSuite) TestUpdateOmitEmpty() {
	idx := New(s.mockClient, &Config{Name: "test"})

	// Note whitespace here! This is NDJSON
	request := []byte(`{"update":{"_index":"test","_id":"objId"}}
{"doc":{}}
`)
	response := []byte(`{
	   "took": 30,
	   "errors": false,
	   "items": [
	      {
	         "update": {
	            "_index": "test",
	            "_type": "_doc",
	            "_id": "objId",
	            "_version": 1,
	            "result": "updated",
	            "_shards": {
	               "total": 2,
	               "successful": 1,
	               "failed": 0
	            },
	            "status": 200,
	            "_seq_no" : 1,
	            "_primary_term" : 2
	         }
	      }
	   ]
	}`)

	type testType struct {
		Field1 *string `json:"field1,omitempty"`
		Field2 []int   `json:"field2,omitempty"`
	}

	dst := testType{
		Field1: nil,
		Field2: []int{},
	}

	testURL := "/_bulk"
	s.MockAPIHandler.
		On("Handle", "POST", testURL, request).
		Return(httpmock.Response{
			Body: response,
		}).
		Once()

	err := idx.Update(s.ctx, "objId", &dst)
	s.NoError(err)

	// Ensure flushing
	s.ctxCancel()
	time.Sleep(100 * time.Millisecond)

	s.MockAPIHandler.AssertExpectations(s.T())
}

func (s *IndexTestSuite) TestDelete() {
	idx := New(s.mockClient, &Config{Name: "test"})

	// Note whitespace here! This is NDJSON
	request := []byte(`{"delete":{"_index":"test","_id":"objId"}}
`)
	response := []byte(`{
	   "took": 30,
	   "errors": false,
	   "items": [
	      {
	         "delete": {
	            "_index": "test",
	            "_type": "_doc",
	            "_id": "objId",
	            "_version": 1,
	            "result": "deleted",
	            "_shards": {
	               "total": 2,
	               "successful": 1,
	               "failed": 0
	            },
	            "status": 200,
	            "_seq_no" : 1,
	            "_primary_term" : 2
	         }
	      }
	   ]
	}`)

	testURL := "/_bulk"
	s.MockAPIHandler.
		On("Handle", "POST", testURL, request).
		Return(httpmock.Response{
			Body: response,
		}).
		Once()

	err := idx.Delete(s.ctx, "objId")
	s.NoError(err)

	// Ensure flushing
	s.ctxCancel()
	time.Sleep(100 * time.Millisecond)

	s.MockAPIHandler.AssertExpectations(s.T())
}

func (s *IndexTestSuite) TestGetFound() {
	idx := New(s.mockClient, &Config{Name: "test"})

	type testType struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	dst := testType{}

	s.mockAsyncGetter.On(
		"Get",
		mock.Anything,
		&bulkgetter.GetRequest{Index: "test", DocumentID: "objId", Fields: []string{"field1", "field2"}},
		&dst,
	).Return(bulkgetter.GetResponse{Found: true, Error: nil})

	result, err := idx.Get(s.ctx, "objId", &dst, "field1", "field2")
	s.NoError(err)
	s.True(result)

	s.mockAsyncGetter.AssertExpectations(s.T())
}

func (s *IndexTestSuite) TestGetNotFound() {
	idx := New(s.mockClient, &Config{Name: "test"})

	type testType struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	dst := testType{}

	s.mockAsyncGetter.On(
		"Get",
		mock.Anything,
		&bulkgetter.GetRequest{Index: "test", DocumentID: "objId", Fields: []string{"field1", "field2"}},
		&dst,
	).Return(bulkgetter.GetResponse{Found: false, Error: nil})

	result, err := idx.Get(s.ctx, "objId", &dst, "field1", "field2")
	s.NoError(err)
	s.False(result)

	s.mockAsyncGetter.AssertExpectations(s.T())
}

func TestIndexTestSuite(t *testing.T) {
	suite.Run(t, new(IndexTestSuite))
}
