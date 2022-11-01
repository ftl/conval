package conval

import (
	"log"
	"strings"

	"github.com/ftl/hamradio/callsign"
)

type BandAndMode struct {
	Band ContestBand
	Mode Mode
}

type Counter struct {
	setup    Setup
	exchange []ExchangeDefinition
	rules    Scoring

	callsignsPerBandAndMode map[BandAndMode]map[callsign.Callsign]bool
	multisPerBandAndMode    map[BandAndMode]map[Property]map[string]bool
	scorePerBand            map[ContestBand]BandScore
}

type QSOScore struct {
	Points           int                      `yaml:"points"`
	Multis           int                      `yaml:"multis"`
	MultiValues      map[Property]string      `yaml:"-"`
	MultiBandAndMode map[Property]BandAndMode `yaml:"-"`
	Duplicate        bool                     `yaml:"duplicate"`
}

func (s QSOScore) Equal(other QSOScore) bool {
	return s.Points == other.Points && s.Multis == other.Multis && s.Duplicate == other.Duplicate
}

type BandScore struct {
	Points int `yaml:"points"`
	Multis int `yaml:"multis"`
}

func NewCounter(setup Setup, exchange []ExchangeDefinition, rules Scoring) *Counter {
	return &Counter{
		setup:    setup,
		exchange: exchange,
		rules:    rules,

		callsignsPerBandAndMode: make(map[BandAndMode]map[callsign.Callsign]bool),
		multisPerBandAndMode:    make(map[BandAndMode]map[Property]map[string]bool),
		scorePerBand:            make(map[ContestBand]BandScore),
	}
}

func (c Counter) TotalScore() BandScore {
	return c.scorePerBand[BandAll]
}

func (c Counter) EffectiveExchangeFields(theirContinent Continent, theirCountry DXCCEntity) []ExchangeField {
	return filterExchangeFields(c.exchange, c.setup.MyContinent, c.setup.MyCountry, theirContinent, theirCountry)
}

func (c *Counter) Add(qso QSO) QSOScore {
	result := c.Probe(qso)

	// apply the QSO band rule
	bandAndMode := effectiveBandAndMode(qso.Band, qso.Mode, c.rules.QSOBandRule)

	// update the callsign registry
	var callsigns map[callsign.Callsign]bool
	var callsignsOK bool
	callsigns, callsignsOK = c.callsignsPerBandAndMode[bandAndMode]
	if !callsignsOK {
		callsigns = make(map[callsign.Callsign]bool)
	}
	callsigns[qso.TheirCall] = true
	c.callsignsPerBandAndMode[bandAndMode] = callsigns

	// update the multi registry
	for property, value := range result.MultiValues {
		bandAndMode := result.MultiBandAndMode[property]
		properties, propertiesOK := c.multisPerBandAndMode[bandAndMode]
		if !propertiesOK {
			properties = make(map[Property]map[string]bool)
		}
		values, valuesOK := properties[property]
		if !valuesOK {
			values = make(map[string]bool)
		}
		values[value] = true
		properties[property] = values
		c.multisPerBandAndMode[bandAndMode] = properties
	}

	// update the scores
	if !result.Duplicate {
		totalScore := c.scorePerBand[BandAll]
		totalScore.Points += result.Points
		totalScore.Multis += result.Multis
		c.scorePerBand[BandAll] = totalScore

		scorePerBand := c.scorePerBand[qso.Band]
		scorePerBand.Points += result.Points
		scorePerBand.Multis += result.Multis
		c.scorePerBand[qso.Band] = scorePerBand
	}

	return result
}

func (c Counter) Probe(qso QSO) QSOScore {
	result := QSOScore{
		MultiValues:      make(map[Property]string),
		MultiBandAndMode: make(map[Property]BandAndMode),
	}

	getProperty := func(property Property) string {
		getter, getterOK := PropertyGetters[property]
		if !getterOK {
			log.Printf("no property getter for %s", property)
			return ""
		}
		return getter.GetProperty(qso)
	}

	// find the relevant QSO rules
	qsoRules := filterScoringRules(c.rules.QSORules, true, c.setup.MyContinent, c.setup.MyCountry, qso.TheirContinent, qso.TheirCountry, qso.Band, getProperty)
	// log.Printf("found %d relevant QSO rules", len(qsoRules))
	if len(qsoRules) == 1 {
		result.Points = qsoRules[0].Value
	} else if len(qsoRules) > 1 {
		value := qsoRules[0].Value
		allEqual := true
		for _, rule := range qsoRules {
			if value != rule.Value {
				allEqual = false
				break
			}
		}
		if allEqual {
			result.Points = value
		} else {
			log.Printf("inconsistent QSO rules: %+v", qsoRules)
		}
	}

	// apply the QSO band rule
	bandAndMode := effectiveBandAndMode(qso.Band, qso.Mode, c.rules.QSOBandRule)

	// check the callsign registry for duplicate
	var callsigns map[callsign.Callsign]bool
	var bandOK bool
	callsigns, bandOK = c.callsignsPerBandAndMode[bandAndMode]
	if bandOK {
		_, result.Duplicate = callsigns[qso.TheirCall]
	}

	// find the relevant multi rules
	multiRules := filterScoringRules(c.rules.MultiRules, false, c.setup.MyContinent, c.setup.MyCountry, qso.TheirContinent, qso.TheirCountry, qso.Band, getProperty)
	// log.Printf("found %d relevant multi rules", len(multiRules))
	for _, rule := range multiRules {
		if rule.Property == "" {
			result.Multis += rule.Value
			continue
		}

		// get the property value
		value := getProperty(rule.Property)

		// apply the band rule
		bandAndMode := effectiveBandAndMode(qso.Band, qso.Mode, rule.BandRule)

		// check for duplicate values
		var duplicateValue bool
		properties, propertiesOK := c.multisPerBandAndMode[bandAndMode]
		if propertiesOK {
			values, propertyOK := properties[rule.Property]
			if propertyOK {
				_, duplicateValue = values[value]
			}
		}

		// count the multi if it is new
		if !duplicateValue {
			result.Multis += rule.Value
			result.MultiValues[rule.Property] = value
			result.MultiBandAndMode[rule.Property] = bandAndMode
		}
	}

	if result.Duplicate {
		result.Points = 0
		result.Multis = 0
	}

	return result
}

func effectiveBandAndMode(band ContestBand, mode Mode, rule BandRule) BandAndMode {
	switch rule {
	case Once:
		return BandAndMode{BandAll, ModeALL}
	case OncePerBand:
		return BandAndMode{band, ModeALL}
	case OncePerBandAndMode:
		return BandAndMode{band, mode}
	default:
		return BandAndMode{}
	}
}

type propertyProvider func(property Property) string

func filterScoringRules(rules []ScoringRule, onlyMostRelevant bool, myContinent Continent, myCountry DXCCEntity, theirContinent Continent, theirCountry DXCCEntity, band ContestBand, getProperty propertyProvider) []ScoringRule {
	matchingRules := make([]ScoringRule, 0, len(rules))
	ruleScores := make([]int, 0, len(matchingRules))
	maxRuleScore := 0
	for _, rule := range rules {
		ruleScore := 0

		if myContinent != "" && len(rule.MyContinent) > 0 {
			if !contains(rule.MyContinent, myContinent) {
				// log.Printf("not my continent %s %v", myContinent, rule.MyContinent)
				continue
			}
			ruleScore++
		}
		if myCountry != "" && len(rule.MyCountry) > 0 {
			if !contains(rule.MyCountry, myCountry) {
				// log.Printf("not my country %s %v", myCountry, rule.MyCountry)
				continue
			}
			ruleScore++
		}
		if theirContinent != "" && len(rule.TheirContinent) > 0 {
			if len(rule.TheirContinent) == 1 &&
				((rule.TheirContinent[0] == SameContinent && myContinent == theirContinent) ||
					(rule.TheirContinent[0] == OtherContinent && myContinent != theirContinent)) {
				ruleScore++
			} else if contains(rule.TheirContinent, theirContinent) {
				ruleScore++
			} else {
				// log.Printf("not their continent %s %v", theirContinent, rule.TheirContinent)
				continue
			}
		}
		if theirCountry != "" && len(rule.TheirCountry) > 0 {
			if len(rule.TheirCountry) == 1 &&
				((rule.TheirCountry[0] == SameCountry && myCountry == theirCountry) ||
					(rule.TheirCountry[0] == OtherCountry && myCountry != theirCountry)) {
				ruleScore++
			} else if contains(rule.TheirCountry, theirCountry) {
				ruleScore++
			} else {
				// log.Printf("not their country %s %v", theirCountry, rule.TheirCountry)
				continue
			}
		}
		if band != "" && len(rule.Bands) > 0 {
			if contains(rule.Bands, band) {
				ruleScore++
			} else {
				// log.Printf("not a valid band %s %v", band, rule.Bands)
				continue
			}
		}
		if rule.TheirWorkingCondition != "" {
			value := strings.ToLower(strings.TrimSpace(getProperty(WorkingConditionProperty)))
			if value != rule.TheirWorkingCondition {
				// log.Printf("not their working condition %q %q", value, rule.TheirWorkingCondition)
				continue
			}
			ruleScore++
		}
		if rule.Property != "" {
			value := getProperty(rule.Property)
			if value == "" {
				// log.Printf("empty property %s", rule.Property)
				continue
			}
			ruleScore++
		}

		ruleScore += rule.AdditionalWeight

		matchingRules = append(matchingRules, rule)
		ruleScores = append(ruleScores, ruleScore)
		if maxRuleScore < ruleScore {
			maxRuleScore = ruleScore
		}
	}

	result := make([]ScoringRule, 0, len(matchingRules))
	for i, rule := range matchingRules {
		if ruleScores[i] == maxRuleScore || !onlyMostRelevant {
			result = append(result, rule)
		}
	}
	return result
}

func filterExchangeFields(definitions []ExchangeDefinition, myContinent Continent, myCountry DXCCEntity, theirContinent Continent, theirCountry DXCCEntity) []ExchangeField {
	matchingDefinitions := make([]ExchangeDefinition, 0, len(definitions))
	definitionScores := make([]int, 0, len(matchingDefinitions))
	maxDefinitionScore := 0

	for _, definition := range definitions {
		definitionScore := 0

		if myContinent != "" && len(definition.MyContinent) > 0 {
			if !contains(definition.MyContinent, myContinent) {
				// log.Printf("not my continent %s %v", myContinent, definition.MyContinent)
				continue
			}
			definitionScore++
		}
		if myCountry != "" && len(definition.MyCountry) > 0 {
			if !contains(definition.MyCountry, myCountry) {
				// log.Printf("not my country %s %v", myCountry, definition.MyCountry)
				continue
			}
			definitionScore++
		}
		if theirContinent != "" && len(definition.TheirContinent) > 0 {
			if len(definition.TheirContinent) == 1 &&
				((definition.TheirContinent[0] == SameContinent && myContinent == theirContinent) ||
					(definition.TheirContinent[0] == OtherContinent && myContinent != theirContinent)) {
				definitionScore++
			} else if contains(definition.TheirContinent, theirContinent) {
				definitionScore++
			} else {
				// log.Printf("not their continent %s %v", theirContinent, definition.TheirContinent)
				continue
			}
		}
		if theirCountry != "" && len(definition.TheirCountry) > 0 {
			if len(definition.TheirCountry) == 1 &&
				((definition.TheirCountry[0] == SameCountry && myCountry == theirCountry) ||
					(definition.TheirCountry[0] == OtherCountry && myCountry != theirCountry)) {
				definitionScore++
			} else if contains(definition.TheirCountry, theirCountry) {
				definitionScore++
			} else {
				// log.Printf("not their country %s %v", theirCountry, definition.TheirCountry)
				continue
			}
		}

		definitionScore += definition.AdditionalWeight

		matchingDefinitions = append(matchingDefinitions, definition)
		definitionScores = append(definitionScores, definitionScore)
		if maxDefinitionScore < definitionScore {
			maxDefinitionScore = definitionScore
		}
	}

	for i, definition := range matchingDefinitions {
		if definitionScores[i] == maxDefinitionScore {
			return definition.Fields
		}
	}
	return nil
}

func contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
