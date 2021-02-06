package config

/*
Workers contains the configuration for the worker pool.

It is fully contained here in order to avoid cyclic imports as the worker package uses the central Config struct.
*/
type Workers struct {
	HashWorkers      int `yaml:"hash_workers" env:"HASH_WORKERS"`
	FileWorkers      int `yaml:"file_workers" env:"FILE_WORKERS"`
	DirectoryWorkers int `yaml:"directory_workers" env:"DIRECTORY_WORKERS"`
}

// WorkersDefaults returns the default configuration for the workerpool.
func WorkersDefaults() Workers {
	return Workers{
		HashWorkers:      70,
		FileWorkers:      120,
		DirectoryWorkers: 70,
	}
}
