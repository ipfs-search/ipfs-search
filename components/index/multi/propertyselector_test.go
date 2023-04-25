package multi

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func makeMetadataProp(mimetype string) Properties {
	return map[string]interface{}{
		"metadata": map[string]interface{}{
			"Content-Type": mimetype,
		},
	}
}

type PropertySelectorTest struct {
	suite.Suite

	ps *PropertySelector
}

func (s *PropertySelectorTest) SetupTest() {
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

func (s *PropertySelectorTest) TestListIndexes() {
	expected := []string{"default", "prefix_matcher", "empty_match", "never_matcher", "multi_matcher"}
	s.Equal(expected, s.ps.ListIndexes())
}

func (s *PropertySelectorTest) TestPrefixMatch() {
	s.Equal("prefix_matcher", s.ps.GetIndex("", makeMetadataProp("prefix/something")))
}

func (s *PropertySelectorTest) TestEmptyMatch() {
	s.Equal("empty_match", s.ps.GetIndex("", makeMetadataProp("")))
}

func (s *PropertySelectorTest) TestDefaultMatch() {
	s.Equal("default", s.ps.GetIndex("", makeMetadataProp("something/else")))
}

func (s *PropertySelectorTest) TestMultiMatch() {
	s.Equal("multi_matcher", s.ps.GetIndex("", makeMetadataProp("ape")))
	s.Equal("multi_matcher", s.ps.GetIndex("", makeMetadataProp("monkey")))
	s.Equal("multi_matcher", s.ps.GetIndex("", makeMetadataProp("primate")))
	s.Equal("multi_matcher", s.ps.GetIndex("", makeMetadataProp("human")))
}

func (s *PropertySelectorTest) TestGetPropValMissing() {
	// When property is missing, getPropVal should return "".
	properties := map[string]interface{}{
		"metadata": map[string]interface{}{},
	}
	s.Empty(s.ps.getPropVal(properties))
}

func (s *PropertySelectorTest) TestGetPropValMissing2() {
	// When not even metadata is defined, getPropVal should also return "".
	properties := map[string]interface{}{}
	s.Empty(s.ps.getPropVal(properties))
}

func TestPropertySelectorTest(t *testing.T) {
	suite.Run(t, new(PropertySelectorTest))
}
