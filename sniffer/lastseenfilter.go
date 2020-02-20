package sniffer

import (
	"log"
	"time"
)

type lastSeenFilter struct {
	resources  map[Resource]time.Time
	Expiration time.Duration
	PruneLen   int
}

func NewLastSeenFilter(expiration time.Duration, pruneLen int) *lastSeenFilter {
	r := make(map[Resource]time.Time)

	return &lastSeenFilter{
		Expiration: expiration,
		PruneLen:   pruneLen,
		resources:  r,
	}
}

func (f *lastSeenFilter) prune() {
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

func (f *lastSeenFilter) Filter(p Provider) bool {
	f.prune()

	lastSeen, present := f.resources[*(p.Resource)]

	if !present {
		// Not present, add it!
		log.Printf("Adding: %v", p)
		f.resources[*(p.Resource)] = p.Date

		return true
	}

	if p.Date.Sub(lastSeen) > f.Expiration {
		// Last seen longer than expiration ago, index it!
		log.Printf("Updating: %v", p)

		f.resources[*(p.Resource)] = p.Date
		return true
	}

	// Too recent, don't index
	log.Printf("Disgarding new: %v - last seen: %s", p, lastSeen)
	return false
}
