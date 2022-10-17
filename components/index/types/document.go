package types

import (
	"time"
)

// Document represents a common properties of resources in an Index.
type Document struct {
	FirstSeen  time.Time  `json:"first-seen"`
	LastSeen   time.Time  `json:"last-seen"`
	References References `json:"references"`
	Size       uint64     `json:"size"`
}
