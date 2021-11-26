package providerfilters

import (
	"log"
	"time"

	t "github.com/ipfs-search/ipfs-search/types"
)

const logEvery = 1000

// LastSeenFilter filters out recently seen Providers.
type LastSeenFilter struct {
	resources  map[string]time.Time
	icount     uint // Iteration counter.
	Expiration time.Duration
	PruneLen   int
}

// NewLastSeenFilter initialises a new LastSeenFilter and returns a pointer to it.
func NewLastSeenFilter(expiration time.Duration, pruneLen int) *LastSeenFilter {
	// Allocate memory for pruneLen+1
	r := make(map[string]time.Time, pruneLen+1)

	return &LastSeenFilter{
		Expiration: expiration,
		PruneLen:   pruneLen,
		resources:  r,
	}
}

func (f *LastSeenFilter) prune() {
	if len(f.resources) > f.PruneLen {
		// Delete all expired items
		now := time.Now()
		cnt := 0

		for i, t := range f.resources {
			if now.Sub(t) > f.Expiration {
				delete(f.resources, i)
				cnt++
			}
		}

		log.Printf("Pruned %d resources, len: %d, pruneLen: %d", cnt, len(f.resources), f.PruneLen)
	}
}

func (f *LastSeenFilter) shouldLog() bool {
	return f.icount%logEvery == 0
}

// Filter takes a Provider and returns true when it is to be included, false
// when not and an error when unexpected condition occur.
func (f *LastSeenFilter) Filter(p t.Provider) (bool, error) {
	f.icount++

	f.prune()

	lastSeen, present := f.resources[p.Resource.String()]

	if !present {
		// Not present, add it!
		if f.shouldLog() {
			log.Printf("Adding LastSeen: %v, len: %d", p, len(f.resources))
		}
		f.resources[p.Resource.String()] = p.Date

		// Index it!
		return true, nil
	}

	if p.Date.Sub(lastSeen) > f.Expiration {
		// Last seen longer than expiration ago, update last seen.
		if f.shouldLog() {
			log.Printf("Updating LastSeen: %v, len: %d", p, len(f.resources))
		}

		f.resources[p.Resource.String()] = p.Date

		// Index it!
		return true, nil
	}

	// Too recent, don't index
	if f.shouldLog() {
		log.Printf("Filtering recent %v, LastSeen %s", p, lastSeen)
	}
	return false, nil
}
