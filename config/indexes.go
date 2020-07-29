package config

import (
    "github.com/ipfs-search/ipfs-search/index/elasticsearch"
)

// Indexes represents the various indexes we're using
type Indexes map[string]*elasticsearch.Config

// IndexesDefaults returns the default indexes.
func IndexesDefaults() Indexes {
    return Indexes{
        "files": &elasticsearch.Config{
            Name: "ipfs_files",
        },
        "directories": &elasticsearch.Config{
            Name: "ipfs_directories",
        },
        "invalids": &elasticsearch.Config{
            Name: "ipfs_invalids",
        },
    }
}
