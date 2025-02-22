package key

import (
	"fmt"
	"github.com/Ackites/KillWxapkg/api"
	"regexp"
	"strings"
	"sync"
)

var (
	rulesInstance *Rules
	once          sync.Once
	jsonMutex     sync.Mutex
)

func getRulesInstance() (*Rules, error) {
	var err error
	once.Do(func() {
		rulesInstance, err = ReadRuleFile()
	})
	return rulesInstance, err
}

func MatchRules(input string) error {
	rules, err := getRulesInstance()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, rule := range rules.Rules {
		if rule.Enabled {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				return fmt.Errorf("failed to compile regex for rule %s: %v", rule.Id, err)
			}
			matches := re.FindAllStringSubmatch(input, -1)
			for _, match := range matches {
				if len(match) > 0 {
					if strings.TrimSpace(match[0]) == "" {
						continue
					}
					err := api.LogHtml(rule.Id + ":" + match[0])
					if err != nil {
						return fmt.Errorf("failed to append to JSON: %v", err)
					}
				}
			}
		}
	}

	return nil
}
