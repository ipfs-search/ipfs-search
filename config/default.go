package config

// Default returns default configuration.
func Default() *Config {
	return &Config{
		IPFSDefaults(),
		OpenSearchDefaults(),
		AMQPDefaults(),
		TikaDefaults(),
		NSFWDefaults(),
		InstrDefaults(),
		CrawlerDefaults(),
		SnifferDefaults(),
		IndexesDefaults(),
		QueuesDefaults(),
		WorkersDefaults(),
	}
}
