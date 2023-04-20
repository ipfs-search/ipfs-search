package factory

import (
	"net/http"
	"testing"

	"github.com/dankinder/httpmock"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var indexes = []string{
	"archives",
	"audio",
	"data",
	"directories",
	"documents",
	"images",
	"invalids",
	"links",
	"other",
	"partials",
	"unknown",
	"videos",
}

type FactoryTestSuite struct {
	suite.Suite

	// Mock search
	mockAPIHandler *httpmock.MockHandler
	mockAPIServer  *httpmock.Server
	responseHeader http.Header

	f *Factory
}

func (s *FactoryTestSuite) expectHelloWorld() {
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

func (s *FactoryTestSuite) SetupTest() {
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

	s.f = New(client)
}

func (s *FactoryTestSuite) TestGetDesiredMapping() {
	for _, index := range indexes {
		mapping, err := s.f.getDesiredMapping(index)
		s.NoError(err)
		s.NotNil(mapping)
	}
}

func TestFactoryTestSuite(t *testing.T) {
	suite.Run(t, new(FactoryTestSuite))
}
