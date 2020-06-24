package config

import (
    "encoding/json"
    "github.com/ipfs-search/ipfs-search/index/elasticsearch"
)

// Indexes represents the various indexes we're using
type Indexes map[string]*elasticsearch.Config

// IndexesDefaults returns the default indexes.
func IndexesDefaults() Indexes {
    var indexSettings, fileMapping, dirMapping, invalidMapping map[string]interface{}

    if err := json.Unmarshal([]byte(indexSettingsJSON), &indexSettings); err != nil {
        panic(err)
    }
    if err := json.Unmarshal([]byte(fileMappingJSON), &fileMapping); err != nil {
        panic(err)
    }
    if err := json.Unmarshal([]byte(dirMappingJSON), &dirMapping); err != nil {
        panic(err)
    }
    if err := json.Unmarshal([]byte(invalidMappingJSON), &invalidMapping); err != nil {
        panic(err)
    }

    return Indexes{
        "files": &elasticsearch.Config{
            Name:     "ipfs_files_v0",
            Settings: indexSettings,
            Mapping:  fileMapping,
        },
        "directories": &elasticsearch.Config{
            Name:     "ipfs_directories_v0",
            Settings: indexSettings,
            Mapping:  dirMapping,
        },
        "invalids": &elasticsearch.Config{
            Name:     "ipfs_invalids_v0",
            Settings: indexSettings,
            Mapping:  invalidMapping,
        },
    }
}
