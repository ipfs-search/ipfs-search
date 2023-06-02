package testsuite

import (
	"net/http"

	"github.com/dankinder/httpmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Suite is a mixin for mocking OpenSearch.
type Suite struct {
	suite.Suite

	MockAPIHandler *httpmock.MockHandler
	MockAPIServer  *httpmock.Server
	ResponseHeader http.Header
}

func (s *Suite) expectHelloWorld() {
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
	s.MockAPIHandler.
		On("Handle", "GET", "/", mock.Anything).
		Return(httpmock.Response{
			Body: testJSON,
		}).
		Once()

}

func (s *Suite) SetupSearchMock() {
	s.MockAPIHandler = &httpmock.MockHandler{}
	s.MockAPIServer = httpmock.NewServer(s.MockAPIHandler)
	s.ResponseHeader = http.Header{
		"Content-Type": []string{"application/json"},
	}

	s.expectHelloWorld()
}

func (s *Suite) TeardownSearchMock() {
	s.MockAPIServer.Close()
}
