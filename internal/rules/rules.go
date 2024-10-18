package rules

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Modification struct {
	Type  string `yaml:"type"`
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type Rule struct {
	URLPattern    string                 `yaml:"url_pattern"`
	Action        string                 `yaml:"action"`
	Modifications []Modification         `yaml:"modifications,omitempty"`
	ActionValue   map[string]interface{} `yaml:"action_value,omitempty"`
}

type RuleList struct {
	Rules []Rule `yaml:"rules"`
}

func LoadRules(filename string) (*RuleList, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var rulesCfg RuleList
	err = yaml.Unmarshal(file, &rulesCfg)
	if err != nil {
		return nil, err
	}

	return &rulesCfg, nil
}
