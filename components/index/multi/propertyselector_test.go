package multi

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PropertySelectorTestSuite struct {
	suite.Suite

	ps *PropertySelector
}

func (s *PropertySelectorTestSuite) SetupTest() {
	s.ps = &PropertySelector{
		PropertyNameParts: []string{"metadata", "Content-Type"},
		Matchers: []*RegexpMatcher{
			NewRegexpMatcher("prefix_matcher", []string{"^prefix"}),
			NewRegexpMatcher("empty_match", []string{"^$"}),
			NewRegexpMatcher("never_matcher", []string{"^prefix"}),
			NewRegexpMatcher("multi_matcher", []string{"^ape", "^monkey", "^primate", "^human"}),
		},
		DefaultIndex: "default",
	}
}

func (s *PropertySelectorTestSuite) TestListIndexes() {
	expected := []string{"default", "prefix_matcher", "empty_match", "never_matcher", "multi_matcher"}
	s.Equal(expected, s.ps.ListIndexes())
}

func (s *PropertySelectorTestSuite) TestPrefixMatch() {
	s.Equal("prefix_matcher", s.ps.GetIndex("", makeMetadataProp("prefix/something")))
}

func (s *PropertySelectorTestSuite) TestEmptyMatch() {
	s.Equal("empty_match", s.ps.GetIndex("", makeMetadataProp("")))
}

func (s *PropertySelectorTestSuite) TestDefaultMatch() {
	s.Equal("default", s.ps.GetIndex("", makeMetadataProp("something/else")))
}

func (s *PropertySelectorTestSuite) TestMultiMatch() {
	s.Equal("multi_matcher", s.ps.GetIndex("", makeMetadataProp("ape")))
	s.Equal("multi_matcher", s.ps.GetIndex("", makeMetadataProp("monkey")))
	s.Equal("multi_matcher", s.ps.GetIndex("", makeMetadataProp("primate")))
	s.Equal("multi_matcher", s.ps.GetIndex("", makeMetadataProp("human")))
}

func (s *PropertySelectorTestSuite) TestGetPropValMissing() {
	// When property is missing, getPropVal should return "".
	properties := map[string]interface{}{
		"metadata": map[string]interface{}{},
	}
	s.Empty(s.ps.getPropVal(properties))
}

func (s *PropertySelectorTestSuite) TestGetPropValMissing2() {
	// When not even metadata is defined, getPropVal should also return "".
	properties := map[string]interface{}{}
	s.Empty(s.ps.getPropVal(properties))
}

func TestPropertySelectorTestSuite(t *testing.T) {
	suite.Run(t, new(PropertySelectorTestSuite))
}
