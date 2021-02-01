package config

// Default() returns default configuration
func Default() *Config {
    return &Config{
        IPFSDefaults(),
        ElasticSearch{
            URL: "http://localhost:9200",
        },
        AMQPDefaults(),
        TikaDefaults(),
        InstrDefaults(),
        CrawlerDefaults(),
        SnifferDefaults(),
        IndexesDefaults(),
        QueuesDefaults(),
        WorkersDefaults(),
    }
}
