package factory

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ipfs-search/ipfs-search/components/index/opensearch/client"
	"github.com/ipfs-search/ipfs-search/components/index/opensearch/testsuite"
	"github.com/ipfs-search/ipfs-search/instr"
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

	clientConfig := &client.Config{
		URL:   s.MockAPIServer.URL(),
		Debug: true,
	}
	instr := instr.New()
	client, _ := client.New(clientConfig, instr)

	s.f = New(client)
}

func (s *FactoryTestSuite) TeardownTest() {
	s.TeardownSearchMock()
}

func (s *FactoryTestSuite) TestGetDesiredMapping() {
	for _, index := range indexes {
		mapping, err := s.f.getDesiredMapping(index)
		s.NoError(err)
		s.NotNil(mapping)
	}
}

func (s *FactoryTestSuite) TestResolveAlias() {
	for _, index := range indexes {
		mapping, err := s.f.getDesiredMapping(index)
		s.NoError(err)
		s.NotNil(mapping)
	}
}

func TestFactoryTestSuite(t *testing.T) {
	suite.Run(t, new(FactoryTestSuite))
}
