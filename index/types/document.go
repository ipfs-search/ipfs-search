package types

import (
	"time"
)

type Reference struct {
	ParentHash string `json:"parent_hash"`
	Name       string `json:"name"`
}

type References []Reference

type Document struct {
	FirstSeen  time.Time  `json:"first-seen"`
	LastSeen   time.Time  `json:"last-seen"`
	References References `json:"references"`
	Size       uint64     `json:"size"`
}
