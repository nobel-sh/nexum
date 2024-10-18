package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

// server configuration.
type Config struct {
	LogFile    string `yaml:"log_file"`
	ListenAddr string `yaml:"listen_addr"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
