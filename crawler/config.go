package crawler

import (
	"time"
)

// Config contains user configurable options for a crawler
type Config struct {
	RetryWait time.Duration // wait time between retries of failed requests

	MetadataMaxSize uint64 // Don't attempt to get metadata for files over this size

	PartialSize uint64 // Size for partial items - this is the default chunker block size
	// TODO: replace by a sane method of skipping partials
}
