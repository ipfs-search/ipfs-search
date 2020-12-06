package ipfs

// Config specifies the configuration for the IPFS protocol.
type Config struct {
	IPFSAPIURL     string // URL of an IPFS API endpoint (for Ls and Stat calls).
	IPFSGatewayURL string // URL of an IPFS Gateway (to request content).
	// LsReadTimeout  time.Duration // Timeout while waiting between yielded directory entries.
	// StatTimeout    time.Duration // Timeout for Stat call.
}

// DefaultConfig returns the default configuration for a Sniffer.
func DefaultConfig() *Config {
	return &Config{
		IPFSAPIURL:     "http://localhost:5001",
		IPFSGatewayURL: "http://localhost:8080",
	}
}
