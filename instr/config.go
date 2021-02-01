package instr

type Config struct {
	SamplingRatio  float64 // Parent-based sampling ratio (fraction of sniffed hashes traced).
	JaegerEndpoint string  // Send spans to Jaeger HTTP endpoint.
}

func DefaultConfig() *Config {
	return &Config{
		SamplingRatio:  1.0,
		JaegerEndpoint: "http://localhost:14268/api/traces",
	}
}
