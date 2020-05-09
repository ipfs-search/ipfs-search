package sniffer

import (
	"context"
	t "github.com/ipfs-search/ipfs-search/types"
	"time"
)

// Extractor takes a message and returns a Provider when available, or nil, or
// an error, when unexpected data was encountered in the message.
type Extractor interface {
	Extract(message map[string]interface{}) (*t.Provider, error)
}

// Queue allows publishing of sniffed items.
type Queue interface {
	Publish(interface{}, uint8) error
}

// Logger yields log messages to extract messages from.
type Logger interface {
	Next() (map[string]interface{}, error)
	Close() error
}

// Shell allows us to get a Logger and set the timeout, it's implemented by shell.Logger.
type Shell interface {
	SetTimeout(time.Duration)
	// Note: it's rather weird but Golang doesn't accept the Logger interface from
	// above here.
	GetLogs(context.Context) (Logger, error)
}
