package bulkgetter

import (
	"time"

	"github.com/opensearch-project/opensearch-go"
)

// Config provides configuration for a BatchingGetter.
type Config struct {
	Client       *opensearch.Client
	BatchSize    int
	BatchTimeout time.Duration
}
