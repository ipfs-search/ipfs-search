package crawler

import (
	"context"
	"github.com/stretchr/testify/suite"

	"github.com/ipfs-search/ipfs-search/extractor"
	"github.com/ipfs-search/ipfs-search/index"
	"github.com/ipfs-search/ipfs-search/protocol"
	t "github.com/ipfs-search/ipfs-search/types"
)

type CrawlerTestSuite struct {
	suite.Suite

	ctx       context.Context
	indexes   Indexes
	protocol  protocol.Protocol
	extractor extractor.Extractor
	c         *Crawler
}

func (s *CrawlerTestSuite) SetupTest() {
	s.ctx = context.Background()

	// Creat a crawler with mocked dependencies
	s.indexes = Indexes{
		Files:       &index.Mock{},
		Directories: &index.Mock{},
		Unsupported: &index.Mock{},
		Invalid:     &index.Mock{},
	}
	s.protocol = &protocol.Mock{}
	s.extractor = &extractor.Mock{}

	s.c = New(s.indexes, s.protocol, s.extractor)
}

func (s *CrawlerTestSuite) TestCrawlHash() {
	r := *t.AnnotatedResource{
		&Resource{
			Protocol: t.IPFSProtocol,
			ID:       "",
		},
		Reference{
			Parent: &Resource{
				Protocol: t.IPFSProtocol,
				ID:       "",
			},
			Name: "banana",
		},
		Stat{
			Type: t.UndefinedType,
			Size: 0,
		},
	}

}

func (s *CrawlerTestSuite) TestCrawlDirectory() {

}

func (s *CrawlerTestSuite) TestCrawlFile() {

}
