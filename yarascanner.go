package main

import (
	// standard
	"fmt"
	"github.com/hillu/go-yara/v4"
)

type YaraScanner struct {
	rulesets []*RuleSet
}

// To implement an interface in Go all you need to do is just define all the functions in the interface
func (self *YaraScanner) Scan(data [] byte) (*ScanReport, error) {
	var matches []ScanMatch
	response := new(ScanReport)
	response.Filename = "stream"
	response.Matches = matches

	for _, ruleset := range self.rulesets {
		var m yara.MatchRules
		err := ruleset.Rules.ScanMem(data, 0, 300, &m)
		if err != nil {
			response.Status = "ERROR"
			return response,err
		}
		for _, resp := range m {
			var match ScanMatch
			match.Rule = resp.Rule
			match.Namespace = resp.Namespace
			match.Tags = resp.Tags
			matches = append(matches, match)
		}

	}
	if len(matches) > 0 {
		response.Status = "INFECTED"
	} else {
		response.Status = "CLEAN"
		matches = [] ScanMatch{}
	}
	
	response.Matches = matches
	return response,nil
}

func (self *YaraScanner) LoadIndex(indexPath string) error {
	ruleset, err := NewRuleSet(indexPath)
	if err != nil {
		return err
	}
	self.rulesets = append(self.rulesets, ruleset)
	return nil
}

func (self *YaraScanner) ListRuleSets() (*RuleSetResponseObject, error) {
	response := new(RuleSetResponseObject)
	for _, ruleset := range self.rulesets {
		response.Names = append(response.Names, ruleset.Name)
	}

	return response,nil

}

func (self *YaraScanner) ListRules(rulesetname string) (*RuleListResponseObject, error) {
	response := new(RuleListResponseObject)
	fmt.Printf("listRules called, %s\n", rulesetname)
	for _, ruleset := range self.rulesets {
		fmt.Printf("Looking for %s, looking at %s\n", rulesetname, ruleset.Name)
		if ruleset.Name == rulesetname {
			rules, err := ruleset.ListRules()
			if err != nil {
				return nil, err
			}
			for _, rule := range rules {
				response.Rules = append(response.Rules, rule)
			}
		}
	}

	return response, nil
}