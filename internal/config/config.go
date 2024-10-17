package config

import (
	"nexum/internal/rules"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Rules []rules.Rule `yaml:"rules"`
}

func Load(filename string) (*Config, error) {
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
