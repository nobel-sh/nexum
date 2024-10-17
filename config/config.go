package config

import (
	"os"

	"gopkg.in/yaml.v2"
	"nexum/rules"
)

type Config struct {
	Rules []rules.Rule `yaml:"rules"`
}

var config Config

func LoadConfig(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return err
	}

	return nil
}

func GetRules() []rules.Rule {
	return config.Rules
}
