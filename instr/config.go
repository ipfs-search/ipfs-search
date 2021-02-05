package instr

// Config specifies the configuration for the instrumentation.
type Config struct {
	SamplingRatio  float64 // Parent-based sampling ratio (fraction of sniffed hashes traced).
	JaegerEndpoint string  // Send spans to Jaeger HTTP endpoint.
}

// DefaultConfig returns the default configuration for the instrumentation.
func DefaultConfig() *Config {
	return &Config{
		SamplingRatio:  0.01,
		JaegerEndpoint: "http://localhost:14268/api/traces",
	}
}
