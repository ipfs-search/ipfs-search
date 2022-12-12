package config

import "time"

/*
Workers contains the configuration for the worker pool.

It is fully contained here in order to avoid cyclic imports as the worker package uses the central Config struct.
*/
type Workers struct {
	HashWorkers       int           `yaml:"hash_workers" env:"HASH_WORKERS"`
	FileWorkers       int           `yaml:"file_workers" env:"FILE_WORKERS"`
	DirectoryWorkers  int           `yaml:"directory_workers" env:"DIRECTORY_WORKERS"`
	MaxIPFSConns      int           `yaml:"ipfs_max_connections" env:"IPFS_MAX_CONNECTIONS"`
	MaxExtractorConns int           `yaml:"extractor_max_connections" env:"EXTRACTOR_MAX_CONNECTIONS"`
	MaxLoadRatio      float64       `yaml:"throttle_max_load" env:"THROTTLE_MAX_LOAD"` // Maximum system load before throttling / load limiting kicks in, as 1-minute load divided per CPU.
	ThrottleMin       time.Duration `yaml:"throttle_min_wait" env:"THROTTLE_MIN_WAIT"` // Minimum time to wait when throttling. The actual time doubles whenever load is above max until reaching ThrottleMax.
	ThrottleMax       time.Duration `yaml:"throttle_max_wait" env:"THROTTLE_MAX_WAIT"` // Maximum time to wait when throttling.
}

// WorkersDefaults returns the default configuration for the workerpool.
func WorkersDefaults() Workers {
	return Workers{
		HashWorkers:       70,
		FileWorkers:       120,
		DirectoryWorkers:  70,
		MaxIPFSConns:      1000,
		MaxExtractorConns: 100,
		MaxLoadRatio:      0.7,
		ThrottleMin:       2 * time.Second,
		ThrottleMax:       time.Minute,
	}
}
