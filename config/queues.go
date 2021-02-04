package config

// Queue holds the configuration for a single Queue.
type Queue struct {
	Name string `yaml:"name` // Name of the Queue.
}

// Queues represents the various queues we're using
type Queues struct {
	Files       Queue `yaml:"files"`       // Resources known to be files.
	Directories Queue `yaml:"directories"` // Resources known to be directories.
	Hashes      Queue `yaml:"hashes"`      // Resources with unknown type.
}

// QueuesDefaults returns the default queues.
func QueuesDefaults() Queues {
	return Queues{
		Files: Queue{
			Name: "files",
		},
		Directories: Queue{
			Name: "directories",
		},
		Hashes: Queue{
			Name: "hashes",
		},
	}
}
