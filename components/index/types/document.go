package types

import (
	"time"
)

// Reference represents a named reference to a Document.
type Reference struct {
	ParentHash string `json:"parent_hash"`
	Name       string `json:"name"`
}

// References is a collection of references to a Document.
type References []Reference

// Document represents a common properties of resources in an Index.
type Document struct {
	FirstSeen  time.Time  `json:"first-seen"`
	LastSeen   time.Time  `json:"last-seen"`
	References References `json:"references"`
	Size       uint64     `json:"size"`
}
