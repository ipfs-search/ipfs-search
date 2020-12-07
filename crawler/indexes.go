package crawler

import (
	"github.com/ipfs-search/ipfs-search/index"
)

type Indexes struct {
	Files       index.Index
	Directories index.Index
	Invalids    index.Index
}
