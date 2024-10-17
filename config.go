package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"nexum/proxy"
)

type Config struct {
	Rules []rules.Rule `yaml:"rules"`
}

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

func WriteLogEntry(entry string) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorf("Error opening log file: %v", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(entry); err != nil {
		log.Errorf("Error writing to log file: %v", err)
	}
}
