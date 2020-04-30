package sniffer

import "time"

// Config holds configuration elements.
type Config struct {
	IpfsAPI            string        // IPFS API endpoint
	AMQPURL            string        // AMPQ URL
	LastSeenExpiration time.Duration // Expiration time for the last-seen resources
	LastSeenPruneLen   int           // Cleanup expired resources from the last-seen
	LoggerTimeout      time.Duration // Throw timeout error when no log messages arrive
}
