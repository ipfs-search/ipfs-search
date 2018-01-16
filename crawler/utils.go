package crawler

import (
	"time"
)

// nowISO returns the current date and time in ISO format
func nowISO() string {
	return time.Now().UTC().Format(time.RFC3339)
}
