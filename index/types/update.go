package types

import (
	"time"
)

type Update struct {
	LastSeen   time.Time  `json:"last-seen"`
	References References `json:"references,omitempty"`
}
