package multi

type Properties interface{}

type Selector interface {
	ListIndexes() []string
	GetIndex(id string, properties Properties) string
}
