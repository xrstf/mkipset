package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/xrstf/mkipset/pkg/ip"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Whitelist     []string      `yaml:"whitelist"`
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

	if len(c.SetName) > 20 {
		return errors.New("setNames must be no longer than 20 characters")
	}

	for i, entry := range c.Whitelist {
		if _, err := ip.Parse(entry); err != nil {
			return fmt.Errorf("whitelist item %d is invalid: %v", i+1, err)
		}
	}

	if c.FlushInterval == 0 {
		c.FlushInterval = 60 * time.Second
	}

	if c.FlushInterval < 5*time.Second {
		return fmt.Errorf("flush interval is too short (%s), must be at least 5s", c.FlushInterval)
	}

	return nil
}

func (c *Config) WhitelistIPs() ip.Slice {
	set := ip.NewSet()

	for _, entry := range c.Whitelist {
		parsed, _ := ip.Parse(entry)
		set.Add(*parsed)
	}

	return set.Sorted()
}
