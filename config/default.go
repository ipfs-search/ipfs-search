package config

// Default() returns default configuration
func Default() *Config {
    return &Config{
        IPFSDefaults(),
        ElasticSearch{
            URL: "http://localhost:9200",
        },
        AMQP{
            URL: "amqp://guest:guest@localhost:5672/",
        },
        TikaDefaults(),
        CrawlerDefaults(),
        SnifferDefaults(),
        IndexesDefaults(),
        QueuesDefaults(),
        WorkersDefaults(),
    }
}
