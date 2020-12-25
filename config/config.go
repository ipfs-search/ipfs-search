package config

import (
	"fmt"
	env "github.com/Netflix/go-env"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type ElasticSearch struct {
	URL string `yaml:"url" env:"ELASTICSEARCH_URL"`
}

// Config contains the configuration for commands.
type Config struct {
	IPFS          `yaml:"ipfs"`
	ElasticSearch `yaml:"elasticsearch"`
	AMQP          `yaml:"amqp"`
	Tika          `yaml:"tika"`

	Crawler `yaml:"crawler"`
	Sniffer `yaml:"sniffer"`
	Indexes `yaml:"indexes"`
	Queues  `yaml:"queues"`
	Workers `yaml:"workers"`
}

// String renders config as YAML
// TODO: Consider TextMarshaler
func (c *Config) String() string {
	bs, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("unable to marshal config to YAML: %v", err)
	}
	return string(bs)
}

// ReadFromFile reads configuration options from specified YAML file
// TODO: Consider TextUnmarshaler
func (c *Config) ReadFromFile(filename string) error {
	yamlFile, err := ioutil.ReadFile(filename)
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

	err = ioutil.WriteFile(configFile, bytes, 0644)
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
