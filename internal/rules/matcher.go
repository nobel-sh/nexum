package rules

import (
	"net/http"
	"regexp"
)

func MatchRule(rules RuleList, url string) *Rule {
	for _, rule := range rules.Rules {
		matched, err := regexp.MatchString("^"+rule.URLPattern, url)
		if err != nil {
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
