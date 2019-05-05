package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Protected     []string      `yaml:"protected"`
	FlushInterval time.Duration `yaml:"flushInterval"`
	SetName       string        `yaml:"setName"`
}

func LoadFromFile(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	conf := &Config{}
	if err = yaml.NewDecoder(f).Decode(conf); err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %v", err)
	}

	if err = conf.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %v", err)
	}

	return conf, nil
}

func (c *Config) Validate() error {
	if len(c.SetName) == 0 {
		return errors.New("no ipset setName configured")
	}

	if c.FlushInterval == 0 {
		c.FlushInterval = 60 * time.Second
	}

	if c.FlushInterval < 5*time.Second {
		return fmt.Errorf("flush interval is too short (%s), must be at least 5s", c.FlushInterval)
	}

	return nil
}
