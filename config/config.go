/*
Package config provides central and canonical representation, reading, parsing
and validation of configuration for components.

Configuration consists of 2 representations:

1. The component-specific configuration Config-structs and DefaultConfig()
default generators, which have no dependencies on nor awareness of other
components or their configuration.

2. Central and canonical configuration (this package), importing and wrapping
the various component configuration, wiring it into a single Config struct.

For example, the crawler package contains a Config struct as well as a
DefaultConfig() function. The config package contains a Crawler struct, wrapping
the Config from the crawler package, and CrawlerDefaults() function wrapping the
DefaultConfig() function from the crawler package. The wrapped Crawler struct
provides tags for reading configuration values from a YAML configuration file
and/or OS environment variables.

The Crawler-struct is then contained within the Config-struct, which performs
reading, parsing and validation of configuration files as well as OS environment
variables. In order to acquire the crawler-specific configuration from the Config-struct,
the CrawlerConfig() method must be called.
*/
package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	env "github.com/ipfs-search/go-env"
	yaml "gopkg.in/yaml.v3"
)

// Config contains the configuration for all components.
type Config struct {
	IPFS       `yaml:"ipfs"`
	OpenSearch `yaml:"opensearch"`
	Redis      `yaml:"redis"`
	AMQP       `yaml:"amqp"`
	Tika       `yaml:"tika"`
	NSFW       `yaml:"nsfw"`

	Instr   `yaml:"instrumentation"`
	Crawler `yaml:"crawler"`
	Sniffer `yaml:"sniffer"`
	Indexes `yaml:"indexes"`
	Queues  `yaml:"queues"`
	Workers `yaml:"workers"`
}

// String renders config as YAML
func (c *Config) String() string {
	bs, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("unable to marshal config to YAML: %v", err)
	}
	return string(bs)
}

// ReadFromFile reads configuration options from specified YAML file
func (c *Config) ReadFromFile(filename string) error {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}

	return nil
}

// ReadFromEnv reads configuration options from environment
func (c *Config) ReadFromEnv() error {
	_, err := env.UnmarshalFromEnviron(c)

	if err != nil {
		return err
	}
	return nil
}

// Check configuration file integrity.
func (c *Config) Check() error {
	zeroElements := findZeroElements(*c)
	if len(zeroElements) > 0 {
		return fmt.Errorf("Missing configuration option(s): %s", strings.Join(zeroElements, ", "))

	}

	return nil
}

// Marshall returns the config serialized to bytes[]
func (c *Config) Marshall() ([]byte, error) {
	return yaml.Marshal(c)
}

// Write writes configuration to file as YAML.
func (c *Config) Write(configFile string) error {
	bytes, err := c.Marshall()
	if err != nil {
		return err
	}

	err = os.WriteFile(configFile, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Dump writes configuration to standard output.
func (c *Config) Dump() error {
	bytes, err := c.Marshall()
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(bytes)

	return err
}

// Get configuration from defaults, optional configuration file, or environment.
func Get(configFile string) (*Config, error) {
	// Start with empty configuration
	cfg := Default()

	if configFile != "" {
		fmt.Printf("Reading configuration file: %s\n", configFile)

		err := cfg.ReadFromFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("Error reading configuration file: %v", err)
		}
	}

	// Read configuration values from env
	err := cfg.ReadFromEnv()
	if err != nil {
		return nil, fmt.Errorf("Error reading configuration from env: %v", err)
	}

	return cfg, nil
}
