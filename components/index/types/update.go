package types

import (
	"time"
)

// Update represents the updatable part of a Document.
type Update struct {
	LastSeen   *time.Time `json:"last-seen,omitempty"`
	References References `json:"references,omitempty"`
}
