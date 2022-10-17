package config

// Index represents the configuration for a single Index.
type Index struct {
	Name   string
	Prefix string
}

// Indexes represents the various indexes we're using
type Indexes struct {
	Files       Index `yaml:"files"`
	Directories Index `yaml:"directories"`
	Invalids    Index `yaml:"invalids"`
	Partials    Index `yaml:"partials"`
}

// IndexesDefaults returns the default indexes.
func IndexesDefaults() Indexes {
	return Indexes{
		Files: Index{
			Name:   "ipfs_files",
			Prefix: "f",
		},
		Directories: Index{
			Name:   "ipfs_directories",
			Prefix: "d",
		},
		Invalids: Index{
			Name:   "ipfs_invalids",
			Prefix: "i",
		},
		Partials: Index{
			Name:   "ipfs_partials",
			Prefix: "p",
		},
	}
}
