package conval

import "github.com/ftl/hamradio/callsign"

type Counter struct {
	setup Setup
	rules Scoring

	callsignsPerBand map[ContestBand]map[callsign.Callsign]bool
	multisPerBand    map[ContestBand]map[Property]map[string]bool
	scorePerBand     map[ContestBand]BandScore
}

type QSOScore struct {
	Points      int                      `yaml:"points"`
	Multis      int                      `yaml:"multis"`
	MultiValues map[Property]string      `yaml:"-"`
	MultiBand   map[Property]ContestBand `yaml:"-"`
	Duplicate   bool                     `yaml:"duplicate"`
}

func (s QSOScore) Equal(other QSOScore) bool {
	return s.Points == other.Points && s.Multis == other.Multis && s.Duplicate == other.Duplicate
}

type BandScore struct {
	Points int `yaml:"points"`
	Multis int `yaml:"multis"`
}

func NewCounter(setup Setup, rules Scoring) *Counter {
	return &Counter{
		setup: setup,
		rules: rules,

		callsignsPerBand: make(map[ContestBand]map[callsign.Callsign]bool),
		multisPerBand:    make(map[ContestBand]map[Property]map[string]bool),
		scorePerBand:     make(map[ContestBand]BandScore),
	}
}

func (c Counter) TotalScore() BandScore {
	return c.scorePerBand[BandAll]
}

func (c *Counter) Add(qso QSO) QSOScore {
	result := c.Probe(qso)

	// apply the QSO band rule
	band := effectiveBand(qso.Band, c.rules.QSOBandRule)

	// update the callsign registry
	var callsigns map[callsign.Callsign]bool
	var callsignsOK bool
	callsigns, callsignsOK = c.callsignsPerBand[band]
	if !callsignsOK {
		callsigns = make(map[callsign.Callsign]bool)
	}
	callsigns[qso.TheirCall] = true
	c.callsignsPerBand[band] = callsigns

	// update the multi registry
	for property, value := range result.MultiValues {
		band := result.MultiBand[property]
		properties, propertiesOK := c.multisPerBand[band]
		if !propertiesOK {
			properties = make(map[Property]map[string]bool)
		}
		values, valuesOK := properties[property]
		if !valuesOK {
			values = make(map[string]bool)
		}
		values[value] = true
		properties[property] = values
		c.multisPerBand[band] = properties
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
		MultiValues: make(map[Property]string),
		MultiBand:   make(map[Property]ContestBand),
	}

	getProperty := func(property Property) string {
		getter, getterOK := PropertyGetters[property]
		if !getterOK {
			return ""
		}
		return getter.GetProperty(qso)
	}

	// find the relevant QSO rules
	qsoRules := filterScoringRules(c.rules.QSORules, true, c.setup.MyContinent, c.setup.MyCountry, qso.TheirContinent, qso.TheirCountry, qso.Band, getProperty)
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
		}
	}

	// apply the QSO band rule
	band := effectiveBand(qso.Band, c.rules.QSOBandRule)

	// check the callsign registry for duplicate
	var callsigns map[callsign.Callsign]bool
	var bandOK bool
	callsigns, bandOK = c.callsignsPerBand[band]
	if bandOK {
		_, result.Duplicate = callsigns[qso.TheirCall]
	}

	// find the relevant multi rules
	multiRules := filterScoringRules(c.rules.MultiRules, true, c.setup.MyContinent, c.setup.MyCountry, qso.TheirContinent, qso.TheirCountry, qso.Band, getProperty)
	for _, rule := range multiRules {
		if rule.Property == "" {
			result.Multis += rule.Value
			continue
		}

		// get the property value
		value := getProperty(rule.Property)

		// apply the band rule
		band := effectiveBand(qso.Band, rule.BandRule)

		// check for duplicate values
		var duplicateValue bool
		properties, propertiesOK := c.multisPerBand[band]
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
			result.MultiBand[rule.Property] = band
		}
	}

	if result.Duplicate {
		result.Points = 0
		result.Multis = 0
	}

	return result
}

func effectiveBand(band ContestBand, rule BandRule) ContestBand {
	switch rule {
	case Once:
		return BandAll
	case OncePerBand:
		return band
	default:
		return ""
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
				continue
			}
			ruleScore++
		}
		if myCountry != "" && len(rule.MyCountry) > 0 {
			if !contains(rule.MyCountry, myCountry) {
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
				continue
			}
		}
		if band != "" && len(rule.Bands) > 0 {
			if contains(rule.Bands, band) {
				ruleScore++
			} else {
				continue
			}
		}
		if rule.Property != "" {
			value := getProperty(rule.Property)
			if value == "" {
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

func contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
