package crawler

import (
	"github.com/ipfs-search/ipfs-search/index"
)

// Indexes used for crawling.
type Indexes struct {
	Files       index.Index
	Directories index.Index
	Invalids    index.Index
}
