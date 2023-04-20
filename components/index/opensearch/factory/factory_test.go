package factory

import (
	"testing"

	"github.com/ipfs-search/ipfs-search/components/index/opensearch/testsuite"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
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
	testsuite.Suite

	f *Factory
}

func (s *FactoryTestSuite) SetupTest() {
	s.SetupSearchMock()

	client, _ := opensearch.NewClient(opensearch.Config{
		Addresses: []string{s.MockAPIServer.URL()},
	})

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
