package client

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ipfs-search/ipfs-search/instr"
)

type ClientTestSuite struct {
	suite.Suite

	instr *instr.Instrumentation
	f     *Client
}

func (s *ClientTestSuite) SetupTest() {
	s.instr = instr.New()
}

func (s *ClientTestSuite) TestNew() {
	config := &Config{}
	c, err := New(config, s.instr)
	s.NoError(err)
	s.NotNil(c)
	s.NotNil(c.SearchClient)
	s.NotNil(c.BulkIndexer)
	s.NotNil(c.AliasResolver)
	s.NotNil(c.BulkGetter)
}

func (s *ClientTestSuite) TestSetSearchClient() {
	c := &Client{}
	s.Nil(c.SearchClient)

	cfg := &Config{}
	s.NoError(c.setSearchClient(cfg))
	s.NotNil(c.SearchClient)
}

func (s *ClientTestSuite) TestSetBulkIndexer() {
	c := &Client{}
	s.Nil(c.BulkIndexer)

	cfg := &Config{}
	s.NoError(c.setSearchClient(cfg))
	s.NoError(c.setBulkIndexer(cfg))
	s.NotNil(c.BulkIndexer)
}

func (s *ClientTestSuite) TestSetAliasResolver() {
	c := &Client{}
	s.Nil(c.AliasResolver)

	cfg := &Config{}
	s.NoError(c.setSearchClient(cfg))
	s.NoError(c.setAliasResolver(cfg))
	s.NotNil(c.AliasResolver)
}

func (s *ClientTestSuite) TestSetBulkGetter() {
	c := &Client{}
	s.Nil(c.BulkGetter)

	cfg := &Config{}
	s.NoError(c.setSearchClient(cfg))
	s.NoError(c.setAliasResolver(cfg))
	s.NoError(c.setBulkGetter(cfg))
	s.NotNil(c.BulkGetter)
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
