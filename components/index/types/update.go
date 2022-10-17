package types

import (
	"time"
)

// Update represents the updatable part of a Document.
type Update struct {
	LastSeen   *time.Time `json:"last-seen,omitempty" redis:"l,omitempty"`
	References References `json:"references,omitempty" redis:"r,omitempty"`
}
