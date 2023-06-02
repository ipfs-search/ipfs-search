package multi

import "regexp"

type RegexpMatcher struct {
	Name    string
	regexes []*regexp.Regexp
}

func (p *RegexpMatcher) Match(property string) bool {
	for _, regex := range p.regexes {
		if regex.MatchString(property) {
			return true
		}
	}

	return false
}

func compileRegexes(input []string) []*regexp.Regexp {
	output := make([]*regexp.Regexp, len(input))
	for i, s := range input {
		output[i] = regexp.MustCompile(s)
	}
	return output
}

func NewRegexpMatcher(name string, regexes []string) *RegexpMatcher {
	return &RegexpMatcher{
		Name:    name,
		regexes: compileRegexes(regexes),
	}
}
