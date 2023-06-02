package multi

type PropertySelector struct {
	PropertyNameParts []string
	Matchers          []*RegexpMatcher
	DefaultIndex      string
}

func (m *PropertySelector) getPropVal(props Properties) string {
	innerProps := props

	for _, partName := range m.PropertyNameParts {
		propMap, ok := innerProps.(map[string]interface{})
		if !ok {
			return ""
		}

		innerProps = propMap[partName]
	}

	if result, ok := innerProps.(string); ok {
		return result
	}

	return ""
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
		propVal := m.getPropVal(properties)
		if matcher.Match(propVal) {
			return matcher.Name
		}
	}

	return m.DefaultIndex
}

// Compile-time assurance that implementation satisfies interface.
var _ Selector = &PropertySelector{}
