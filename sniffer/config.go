package sniffer

import "time"

// Config holds configuration elements.
type Config struct {
	IpfsAPI            string
	AMQPURL            string
	LastSeenExpiration time.Duration
	LastSeenPruneLen   int
}
