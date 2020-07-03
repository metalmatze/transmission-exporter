package config

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config represents the configuration for the exporter
type Config struct {
	Clients []*Client `yaml:"clients"`
	WebAddr string    `yaml:"webaddr"`
	WebPath string    `yaml:"webpath"`
}

// Client configration
type Client struct {
	ClientName           string `yaml:"name,omitempty"`
	TransmissionAddr     string `yaml:"addr"`
	TransmissionUsername string `yaml:"username,omitempty"`
	TransmissionPassword string `yaml:"password,omitempty"`
}

// Load reads YAML from reader and unmashals in Config
func Load(r io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	// check name exist, if not, use addr instead
	for _, client := range c.Clients {
		if client.ClientName == "" {
			client.ClientName = client.TransmissionAddr
		}
	}

	return c, nil
}
