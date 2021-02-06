package config

// ElasticSearch holds configuration for ElasticSearch.
type ElasticSearch struct {
	URL string `yaml:"url" env:"ELASTICSEARCH_URL"`
}

// ElasticSearchDefaults returns the defaults for ElasticSearch.
func ElasticSearchDefaults() ElasticSearch {
	return ElasticSearch{
		URL: "http://localhost:9200",
	}
}
