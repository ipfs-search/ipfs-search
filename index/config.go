package index

// Config represents the configuration for a specific index.
type Config struct {
	Name     string
	Settings map[string]interface{}
	Mapping  map[string]interface{}
}

// ConfiguredIndex represents an index with configuration.
type ConfiguredIndex interface {
	GetConfig() *Config
}
