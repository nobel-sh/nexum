package config

import (
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Rule struct {
	URLPattern string `yaml:"url_pattern"`
	Action     string `yaml:"action"` // "allow", "block", or "modify"
}

type Config struct {
	Rules []Rule `yaml:"rules"`
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

func MatchRule(url string) *Rule {
	for _, rule := range config.Rules {
		matched, err := regexp.MatchString(rule.URLPattern, url)
		if err != nil {
			log.Errorf("Error matching rule: %v", err)
			continue
		}
		if matched {
			return &rule
		}
	}
	return nil
}
