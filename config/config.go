package config

import (
	"fmt"
	env "github.com/Netflix/go-env"
	"github.com/c2h5oh/datasize"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/crawler/factory"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Tika struct {
	IpfsTikaURL     string            `yaml:"url" env:"IPFS_TIKA_URL"`
	IpfsTikaTimeout time.Duration     `yaml:"timeout"`
	MetadataMaxSize datasize.ByteSize `yaml:"max_size"`
}

type IPFS struct {
	IpfsAPI     string        `yaml:"api_url" env:"IPFS_API_URL"`
	IpfsTimeout time.Duration `yaml:"timeout"`
}

type ElasticSearch struct {
	ElasticSearchURL string `yaml:"url" env:"ELASTICSEARCH_URL"`
}

type AMQP struct {
	AMQPURL string `yaml:"url" env:"AMQP_URL"`
}

type Crawler struct {
	RetryWait   time.Duration     `yaml:"retry_wait"`
	HashWait    time.Duration     `yaml:"hash_wait"`
	FileWait    time.Duration     `yaml:"file_wait"`
	PartialSize datasize.ByteSize `yaml:"partial_size"`
	HashWorkers uint              `yaml:"hash_workers"`
	FileWorkers uint              `yaml:"file_workers"`
}

type Config struct {
	Tika          `yaml:"tika"`
	IPFS          `yaml:"ipfs"`
	ElasticSearch `yaml:"elasticsearch"`
	AMQP          `yaml:"amqp"`
	Crawler       `yaml:"crawler"`
}

func (c *Config) CrawlerConfig() *crawler.Config {
	return &crawler.Config{
		IpfsTikaURL:     c.Tika.IpfsTikaURL,
		IpfsTikaTimeout: c.Tika.IpfsTikaTimeout,
		MetadataMaxSize: uint64(c.Tika.MetadataMaxSize),
		RetryWait:       c.Crawler.RetryWait,
		PartialSize:     uint64(c.Crawler.PartialSize),
	}
}

func (c *Config) FactoryConfig() *factory.Config {
	return &factory.Config{
		IpfsAPI:          c.IPFS.IpfsAPI,
		IpfsTimeout:      c.IPFS.IpfsTimeout,
		ElasticSearchURL: c.ElasticSearch.ElasticSearchURL,
		AMQPURL:          c.AMQP.AMQPURL,
		CrawlerConfig:    c.CrawlerConfig(),
	}
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

func (c *Config) Check() error {
	zeroElements := findZeroElements(*c)
	if len(zeroElements) > 0 {
		return fmt.Errorf("Missing configuration option(s): %s", strings.Join(zeroElements, ", "))

	}

	return nil
}

func Get(configFile string) (*Config, error) {
	// Start with empty configuration
	cfg := Default()

	err := cfg.ReadFromFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("Error reading configuration file: %v", err)
	}

	// Read configuration values from env
	err = cfg.ReadFromEnv()
	if err != nil {
		return nil, fmt.Errorf("Error reading configuration from env: %v", err)
	}

	// Check configuration before returning
	err = cfg.Check()
	if err != nil {
		return nil, err
	}

	// Print config for debug
	log.Printf("Config:\n%v\n", cfg)

	return cfg, nil
}
