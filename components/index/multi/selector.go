package multi

type Properties map[string]interface{}

type Selector interface {
	ListIndexes() []string
	GetIndex(id string, properties Properties) string
}

type PropertySelector struct {
	PropertyName string
	Matchers     []*RegexpMatcher
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
