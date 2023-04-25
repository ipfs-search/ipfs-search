package multi

import "strings"

type Properties map[string]interface{}

type Selector interface {
	ListIndexes() []string
	GetIndex(id string, properties Properties) string
}

type PrefixMatcher struct {
	Name     string
	Prefixes []string
}

func (p *PrefixMatcher) Match(property string) bool {
	for _, prefix := range p.Prefixes {
		if strings.HasPrefix(property, prefix) {
			return true
		}
	}

	return false
}

type PropertySelector struct {
	PropertyName string
	Matchers     []PrefixMatcher
	DefaultIndex string
}

func (m *PropertySelector) ListIndexes() []string {
	names := make([]string, len(m.Matchers)+1)
	names[0] = m.DefaultIndex
	for i, p := range m.Matchers {
		names[i+1] = p.Name
	}
	return names
}

func (m *PropertySelector) GetIndex(id string, properties Properties) string {
	for _, matcher := range m.Matchers {
		property := properties[m.PropertyName].(string)
		if matcher.Match(property) {
			return matcher.Name
		}
	}

	return m.DefaultIndex
}

// Compile-time assurance that implementation satisfies interface.
var _ Selector = &PropertySelector{}
