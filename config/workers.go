package config

import (
	"time"
)

type Workers struct {
	// RetryWait       time.Duration     `yaml:"retry_wait"` // Time to wait between failed http requests.
	StartupDelay     time.Duration `yaml:"startup_delay"`
	HashWorkers      uint          `yaml:"hash_workers"`
	FileWorkers      uint          `yaml:"file_workers"`
	DirectoryWorkers uint          `yaml:"directory_workers"`
}

func WorkersDefaults() Workers {
	return Workers{
		StartupDelay:     time.Duration(100 * time.Millisecond),
		HashWorkers:      70,
		FileWorkers:      120,
		DirectoryWorkers: 70,
	}
}
