package index

// Config represents the configuration for a specific index.
type Config struct {
	Name     string
	Settings map[string]interface{}
	Mapping  map[string]interface{}
}
