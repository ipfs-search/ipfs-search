package sniffer

type Filter interface {
	Filter(Provider) bool
}
