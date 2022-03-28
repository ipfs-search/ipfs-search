package batchinggetter

import (
	"github.com/opensearch-project/opensearch-go"
)

// Config provides configuration for a BatchingGetter.
type Config struct {
	Client    *opensearch.Client
	BatchSize int
}
