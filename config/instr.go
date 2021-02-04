package config

import (
	"github.com/ipfs-search/ipfs-search/instr"
)

// Instr specifies the configuration for instrumentation.
type Instr struct {
	SamplingRatio  float64 `yaml:"sampling_ratio" env:"OTEL_TRACE_SAMPLER_ARG"`         // Parent-based sampling ratio (fraction of sniffed hashes traced). Defaults to `0.01` (1%). For some reason, setting this as an environment option fails.
	JaegerEndpoint string  `yaml:"jaeger_endpoint" env:"OTEL_EXPORTER_JAEGER_ENDPOINT"` // Send spans to Jaeger HTTP endpoint, for example `http://jaeger:14268/api/traces`.
}

// InstrConfig returns component-specific configuration from the canonical central configuration.
func (c *Config) InstrConfig() *instr.Config {
	cfg := instr.Config(c.Instr)
	return &cfg
}

// InstrDefaults returns the defaults for component configuration, based on the component-specific configuration.
func InstrDefaults() Instr {
	return Instr(*instr.DefaultConfig())
}
