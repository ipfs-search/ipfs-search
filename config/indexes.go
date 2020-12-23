package config

type Index struct {
    Name string
}

// Indexes represents the various indexes we're using
type Indexes struct {
    Files       Index `yaml:"files"`
    Directories Index `yaml:"directories"`
    Invalids    Index `yaml:"invalids"`
}

// IndexesDefaults returns the default indexes.
func IndexesDefaults() Indexes {
    return Indexes{
        Files: Index{
            Name: "ipfs_files",
        },
        Directories: Index{
            Name: "ipfs_directories",
        },
        Invalids: Index{
            Name: "ipfs_invalids",
        },
    }
}
