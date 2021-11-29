package elasticsearch

// TODO: Test whether indexed items with omitempty are actually left out - otherwise
// non-updating references will overwrite the existing!
import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/dankinder/httpmock"
	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type IndexTestSuite struct {
	suite.Suite
	ctx            context.Context
	instr          *instr.Instrumentation
	mockAPIHandler *httpmock.MockHandler
	mockAPIServer  *httpmock.Server
	mockClient     *Client
	responseHeader http.Header
}

func (s *IndexTestSuite) expectHelloWorld() {
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

func (s *IndexTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.instr = instr.New()

	s.mockAPIHandler = &httpmock.MockHandler{}
	s.mockAPIServer = httpmock.NewServer(s.mockAPIHandler)
	s.responseHeader = http.Header{
		"Content-Type": []string{"application/json"},
	}
	config := &ClientConfig{
		URL:   s.mockAPIServer.URL(),
		Debug: true,
	}
	s.mockClient, _ = NewClient(config, s.instr)

	s.expectHelloWorld()
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
	request := []byte(`{"create":{"_id":"objId","_index":"test"}}
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
	s.mockAPIHandler.
		On("Handle", "POST", testURL, request).
		Return(httpmock.Response{
			Body: response,
		}).
		Once()

	err := idx.Index(s.ctx, "objId", &dst)
	s.NoError(err)

	// Ensure flushing
	err = s.mockClient.Close(s.ctx)
	s.NoError(err)

	s.mockAPIHandler.AssertExpectations(s.T())
}

func (s *IndexTestSuite) TestUpdate() {
	idx := New(s.mockClient, &Config{Name: "test"})

	// Note whitespace here! This is NDJSON
	request := []byte(`{"update":{"_id":"objId","_index":"test"}}
{"field1":"hoi","field2":4}

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
	s.mockAPIHandler.
		On("Handle", "POST", testURL, request).
		Return(httpmock.Response{
			Body: response,
		}).
		Once()

	err := idx.Update(s.ctx, "objId", &dst)
	s.NoError(err)

	// Ensure flushing
	err = s.mockClient.Close(s.ctx)
	s.NoError(err)

	s.mockAPIHandler.AssertExpectations(s.T())
}

func (s *IndexTestSuite) TestDelete() {
	idx := New(s.mockClient, &Config{Name: "test"})

	// Note whitespace here! This is NDJSON
	request := []byte(`{"delete":{"_id":"objId","_index":"test"}}
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
	            "status": 202,
	            "_seq_no" : 1,
	            "_primary_term" : 2
	         }
	      }
	   ]
	}`)

	testURL := "/_bulk"
	s.mockAPIHandler.
		On("Handle", "POST", testURL, request).
		Return(httpmock.Response{
			Body: response,
		}).
		Once()

	err := idx.Delete(s.ctx, "objId")
	s.NoError(err)

	// Ensure flushing
	err = s.mockClient.Close(s.ctx)
	s.NoError(err)

	s.mockAPIHandler.AssertExpectations(s.T())
}

func (s *IndexTestSuite) TestGetFound() {
	idx := New(s.mockClient, &Config{Name: "test"})

	testFound := []byte(`{
		"_index": "test",
		"_type": "_doc",
		"_id": "objId",
		"_version": 1,
		"_seq_no": 0,
		"_primary_term": 1,
		"found": true,
		"_source": {
   		"field1": "value",
   		"field2": 5
		}
	}`)

	testURL := "/test/_doc/objId?_source_includes=field1%2Cfield2&preference=_local&realtime=true"
	s.mockAPIHandler.
		On("Handle", "GET", testURL, mock.Anything).
		Return(httpmock.Response{
			Body: testFound,
		}).
		Once()

	type testType struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	dst := testType{}

	result, err := idx.Get(s.ctx, "objId", &dst, "field1", "field2")
	s.NoError(err)
	s.True(result)
	s.Equal(dst, testType{
		Field1: "value",
		Field2: 5,
	})

	s.mockAPIHandler.AssertExpectations(s.T())
}

func (s *IndexTestSuite) TestGetNotFound() {
	idx := New(s.mockClient, &Config{Name: "test"})

	testNotFound := []byte(`{
		"_index": "ipfs_files",
		"_type": "_doc",
		"_id": "objId",
		"found": false
	}`)

	testURL := "/test/_doc/objId?_source_includes=field1%2Cfield2&preference=_local&realtime=true"
	s.mockAPIHandler.
		On("Handle", "GET", testURL, mock.Anything).
		Return(httpmock.Response{
			Body:   testNotFound,
			Status: 404,
		}).
		Once()

	type testType struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	dst := testType{}

	result, err := idx.Get(s.ctx, "objId", &dst, "field1", "field2")
	s.NoError(err)
	s.False(result)
	s.Equal(dst, testType{})

	s.mockAPIHandler.AssertExpectations(s.T())
}

func (s *IndexTestSuite) TestClose() {

}

func TestIndexTestSuite(t *testing.T) {
	suite.Run(t, new(IndexTestSuite))
}
