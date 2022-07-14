package worker

import (
	"time"
)

// Config contains configuration for a Worker.
type Config struct {
	Name         string        // For debugging purposes.
	MaxLoadRatio float64       // Maximum system load before throttling / load limiting kicks in, as 1-minute load divided per CPU.
	ThrottleMin  time.Duration // Minimum time to wait when throttling. The actual time doubles whenever load is above max until reaching ThrottleMax.
	ThrottleMax  time.Duration // Maximum time to wait when throttling.
}
