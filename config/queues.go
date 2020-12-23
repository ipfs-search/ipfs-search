package config

type Queue struct {
	Name string `yaml:"name`
}

type Queues struct {
	Files       Queue `yaml:"files"`
	Directories Queue `yaml:"directories"`
	Hashes      Queue `yaml:"hashes"`
}

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
