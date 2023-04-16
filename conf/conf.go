package conf

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// DefaultPort is used when nothing is specified in config.
const DefaultPort uint = 2023

// C defines a list of configurable options.
type C struct {
	Secret    string   `yaml:"secret"`
	FileStore string   `yaml:"fileStore"`
	Port      uint     `yaml:"port"`
	Reflect   bool     `yaml:"reflect"`
	Endpoints []string `yaml:"endpoints"`
	Tls       bool     `yaml:"tls"`
	Certs     string   `yaml:"certs"`
}

// New unmarshals the given file and returns a new instance of config.
func New(path string) (*C, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read config %q: %w", path, err,
		)
	}

	c := new(C)
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal config %q: %w", path, err,
		)
	}

	return c, nil
}

// Validate validates the config against a set of rules.
func (c *C) Validate() error {
	if c.Secret == "" {
		return fmt.Errorf("missing secret in config")
	}

	if c.FileStore == "" {
		return fmt.Errorf("missing file store path in config")
	}

	if c.Port == 0 {
		c.Port = DefaultPort
	}

	if c.Endpoints == nil || len(c.Endpoints) == 0 {
		return fmt.Errorf("missing endpoints in config")
	}

	if c.Tls && c.Certs == "" {
		return fmt.Errorf("certificate path not set in config")
	}

	return nil
}
