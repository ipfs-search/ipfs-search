package config

import (
	"github.com/c2h5oh/datasize"
)

// ElasticSearch holds configuration for ElasticSearch.
type ElasticSearch struct {
	URL             string            `yaml:"url" env:"ELASTICSEARCH_URL"`
	IndexBufferSize datasize.ByteSize `yaml:"index_buffer_size" env:"ELASTICSEARCH_INDEX_BUFFER_SIZE"`
}

// ElasticSearchDefaults returns the defaults for ElasticSearch.
func ElasticSearchDefaults() ElasticSearch {
	return ElasticSearch{
		URL:             "http://localhost:9200",
		IndexBufferSize: 5 * 1024 * 1024, // 5MB default.
	}
}
