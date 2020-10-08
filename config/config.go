package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Marker struct {
	Mode  string   `yaml:"mode"`
	Marks []uint32 `yaml:"marks"`
}

type Endpoint struct {
	Sniffer string   `yaml:"sniffer"`
	Match   []string `yaml:"match,omitempty"`
	Marker  string   `yaml:"marker"`
}

type Config struct {
	ListenAddress string              `yaml:"listen-address"`
	Markers       map[string]Marker   `yaml:"markers"`
	Endpoints     map[uint16]Endpoint `yaml:"endpoints"`
}

func Parse(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	return cfg, yaml.NewDecoder(file).Decode(cfg)
}
