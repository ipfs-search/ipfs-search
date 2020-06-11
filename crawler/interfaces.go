package crawler

// Queue allows publishing of sniffed items.
type Queue interface {
	Publish(interface{}, uint8) error
}
