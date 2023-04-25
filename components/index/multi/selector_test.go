package multi

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SelectorTestSuite struct {
	suite.Suite
}

func (s *SelectorTestSuite) SetupTest() {
}

func (s *SelectorTestSuite) TestListIndexes() {
}

func TestSelectorTestSuite(t *testing.T) {
	suite.Run(t, new(SelectorTestSuite))
}
