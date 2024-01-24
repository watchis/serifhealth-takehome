package parser

import (
	"encoding/json"
	"errors"
	"os"
	"serifhealth-takehome/config"
	"strings"
)

const REPORTING_STRUCTURE_TOKEN = "reporting_structure"
const KEY_PARAM = "Key-Pair-Id"

type parser struct {
	config *config.Config

	seenPlans map[string]struct{}
}

func NewParser(config *config.Config) *parser {
	return &parser{
		config: config,
		seenPlans: map[string]struct{}{},
	}
}

func (p *parser) ParseFile() ([]string, error) {
	// Setting up the json stream
	reader, err := os.Open(p.config.FilePath)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(reader)

	p.skipHeader(decoder)
	return p.retrieveUrls(decoder)
}

func (p *parser) skipHeader(decoder *json.Decoder) error {
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		
		if token == REPORTING_STRUCTURE_TOKEN {
			// Trimming '[' token from start of reporting structure
			_, err := decoder.Token()
			if err != nil {
				return err
			}

			return nil
		}
	}

	return errors.New("Expected token not found in json header")
}

func (p *parser) retrieveUrls(decoder *json.Decoder) ([]string, error) {
	urlSet := map[string]struct{}{}

	for decoder.More() {
	// for i := 0; i < 100; i++ {
		var report *Report

		err := decoder.Decode(&report)
		if err != nil {
			return nil, err
		}

		if !p.hasUnseenPlan(report.Plans) {
			continue
		}

		for _, file := range report.InNetworkFiles {
			if !p.isTargetState(file.Description) && !p.isAWSDomain(file.Location){
				continue					
			}

			if !p.isValidPlan(file.Description) {
				continue
			}

			urlSet[file.Location] = struct{}{}
		}
	}

	i := 0
	urls := make([]string, len(urlSet))
	for url := range urlSet {
		urls[i] = url
		i++
	}
	return urls, nil
}

func (p *parser) hasUnseenPlan(plans []ReportingPlan) bool {
	for _, plan := range plans {
		if _, ok := p.seenPlans[plan.ID]; ok {
			continue
		}

		p.seenPlans[plan.ID] = struct{}{}
		return true
	}
	
	return false
}

func (p *parser) isAWSDomain(loc string) bool {
	return strings.Contains(loc, "aws") && strings.Contains(loc, p.config.StateAbbreviation)
}

func (p *parser) isTargetState(desc string) bool {
	return strings.Contains(desc, p.config.StateAbbreviation) || strings.Contains(desc, p.config.StateName)
}

func (p *parser) isValidPlan(desc string) bool {
	if p.config.IncludeHighmark && strings.Contains(desc, "Highmark") {
		return true
	}

	return !strings.Contains(desc, "Highmark")
}
