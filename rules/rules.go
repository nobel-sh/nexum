package rules

import (
	"net/http"
	"regexp"

	log "github.com/sirupsen/logrus"
)

type Modification struct {
	Type  string
	Key   string
	Value string
}

type Rule struct {
	URLPattern    string
	Action        string
	Modifications []Modification
}

func MatchRule(rules []Rule, url string) *Rule {
	for _, rule := range rules {
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

func ApplyModifications(r *http.Request, mods []Modification) {
	for _, mod := range mods {
		switch mod.Type {
		case "add_header":
			r.Header.Add(mod.Key, mod.Value)
		case "remove_header":
			r.Header.Del(mod.Key)
		case "set_header":
			r.Header.Set(mod.Key, mod.Value)
		}
	}
}
