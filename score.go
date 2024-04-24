package conval

import (
	"log"
	"sort"
	"strings"

	"github.com/ftl/hamradio/callsign"
)

type BandAndMode struct {
	Band ContestBand
	Mode Mode
}

type Counter struct {
	definition Definition
	setup      Setup
	prefixes   PrefixDatabase
	trace      bool

	qsos                    []ScoredQSO
	callsignsPerBandAndMode map[BandAndMode]map[callsign.Callsign]bool
	multisPerBandAndMode    map[BandAndMode]map[Property]map[string]bool
	scorePerBand            map[ContestBand]BandScore
}

type ScoredQSO struct {
	QSO
	QSOScore
}

type QSOScore struct {
	Points           int                      `yaml:"points" json:"points"`
	Multis           int                      `yaml:"multis" json:"multis"`
	MultiValues      map[Property]string      `yaml:"-" json:"-"`
	MultiBandAndMode map[Property]BandAndMode `yaml:"-" json:"-"`
	Duplicate        bool                     `yaml:"duplicate" json:"duplicate"`
}

func (s QSOScore) Equal(other QSOScore) bool {
	return s.Points == other.Points && s.Multis == other.Multis && s.Duplicate == other.Duplicate
}

type BandScore struct {
	QSOs   int `yaml:"qsos" json:"qsos"`
	Points int `yaml:"points" json:"points"`
	Multis int `yaml:"multis" json:"multis"`
}

func NewCounter(definition Definition, setup Setup, prefixes PrefixDatabase) *Counter {
	return &Counter{
		definition: definition,
		setup:      setup,
		prefixes:   prefixes,

		qsos:                    make([]ScoredQSO, 0, 10000),
		callsignsPerBandAndMode: make(map[BandAndMode]map[callsign.Callsign]bool),
		multisPerBandAndMode:    make(map[BandAndMode]map[Property]map[string]bool),
		scorePerBand:            make(map[ContestBand]BandScore),
	}
}

func (c *Counter) SetTrace(trace bool) {
	c.trace = trace
}

func (c *Counter) tracef(format string, args ...any) {
	if !c.trace {
		return
	}
	log.Printf(format, args...)
}

func (c Counter) UsedBands() []ContestBand {
	result := make([]ContestBand, 0, len(c.scorePerBand)-1)
	for band := range c.scorePerBand {
		if band == BandAll {
			continue
		}
		result = append(result, band)
	}
	return result
}

func (c Counter) MultiProperties() []Property {
	properties := make(map[Property]bool)
	for _, multisPerProperty := range c.multisPerBandAndMode {
		for property := range multisPerProperty {
			properties[property] = true
		}
	}
	result := make([]Property, 0, len(properties))
	for property := range properties {
		result = append(result, property)
	}
	return result
}

func (c Counter) MultisPerBand(property Property, band ContestBand) []string {
	multis := make(map[string]bool)
	for bam, multisPerProperty := range c.multisPerBandAndMode {
		if bam.Band != band {
			continue
		}
		for multi := range multisPerProperty[property] {
			multis[multi] = true
		}
	}

	result := make([]string, 0, len(multis))
	for multi := range multis {
		result = append(result, multi)
	}
	sort.Strings(result)
	return result
}

func (c Counter) MultisPerProperty(property Property) []string {
	multis := make(map[string]bool)
	for _, multisPerProperty := range c.multisPerBandAndMode {
		for multi := range multisPerProperty[property] {
			multis[multi] = true
		}
	}

	result := make([]string, 0, len(multis))
	for multi := range multis {
		result = append(result, multi)
	}
	sort.Strings(result)
	return result
}

func (c Counter) BandsPerMulti(property Property, multi string) []ContestBand {
	bands := make(map[ContestBand]bool)
	for bam, multisPerProperty := range c.multisPerBandAndMode {
		multis := multisPerProperty[property]
		if multis[multi] {
			bands[bam.Band] = true
		}
	}
	result := make([]ContestBand, 0, len(bands))
	for band := range bands {
		result = append(result, band)
	}
	// TODO sort the bands
	return result
}

func (c Counter) BandScore(band ContestBand) BandScore {
	return c.scorePerBand[band]
}

func (c Counter) TotalScore() BandScore {
	return c.scorePerBand[BandAll]
}

func (c Counter) Total(score BandScore) int {
	switch c.definition.Scoring.MultiOperation {
	case AddMultis:
		return score.Points + score.Multis
	default:
		return score.Points * score.Multis
	}
}

func (c *Counter) EffectiveExchangeFields(theirCall callsign.Callsign) []ExchangeField {
	theirContinent, theirCountry, _, _, ok := c.prefixes.Find(theirCall.String())
	c.tracef("filtering exchange fields for %s on continent %s in country %s (found: %t)", theirCall, theirContinent, theirCountry, ok)
	if !ok {
		theirContinent = ""
		theirCountry = ""
	}
	return c.filterExchangeFields(c.definition.Exchange, c.setup.MyContinent, c.setup.MyCountry, theirContinent, theirCountry, theirCall.WorkingCondition)
}

func (c *Counter) Add(qso QSO) QSOScore {
	result := c.Probe(qso)
	c.qsos = append(c.qsos, ScoredQSO{QSO: qso, QSOScore: result})

	// apply the QSO band rule
	bandAndMode := effectiveBandAndMode(qso.Band, qso.Mode, c.definition.Scoring.QSOBandRule)

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
	totalScore := c.scorePerBand[BandAll]
	scorePerBand := c.scorePerBand[qso.Band]

	totalScore.QSOs += 1
	scorePerBand.QSOs += 1
	if !result.Duplicate {
		totalScore.Points += result.Points
		totalScore.Multis += result.Multis

		scorePerBand.Points += result.Points
		scorePerBand.Multis += result.Multis
	}

	c.scorePerBand[BandAll] = totalScore
	c.scorePerBand[qso.Band] = scorePerBand

	return result
}

func (c Counter) Probe(qso QSO) QSOScore {
	c.tracef("probing %+v", qso)

	result := QSOScore{
		MultiValues:      make(map[Property]string),
		MultiBandAndMode: make(map[Property]BandAndMode),
	}

	getMyProperty := func(property Property) string {
		return qso.MyExchange[property]
	}

	getTheirProperty := func(property Property) string {
		getter, getterOK := c.definition.PropertyGetter(property)
		if !getterOK {
			c.tracef("no property getter for %s", property)
			return ""
		}
		return getter.GetProperty(qso, c.prefixes)
	}

	// find the relevant QSO rules
	c.tracef("filtering %d scoring rules", len(c.definition.Scoring.QSORules))
	qsoRules := c.filterScoringRules(c.definition.Scoring.QSORules, true, c.setup.MyContinent, c.setup.MyCountry, c.setup.MyPrefix(), qso.TheirContinent, qso.TheirCountry, qso.TheirPrefix(), qso.Band, getMyProperty, getTheirProperty)
	c.tracef("found %d relevant QSO rules: %+v", len(qsoRules), qsoRules)
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
			c.tracef("inconsistent QSO rules for QSO with %s (%s, %s) at %v: %+v", qso.TheirCall, qso.TheirContinent, qso.TheirCountry, qso.Timestamp, qsoRules)
		}
	}

	// apply the QSO band rule
	bandAndMode := effectiveBandAndMode(qso.Band, qso.Mode, c.definition.Scoring.QSOBandRule)

	// check the callsign registry for duplicate
	var callsigns map[callsign.Callsign]bool
	var bandOK bool
	callsigns, bandOK = c.callsignsPerBandAndMode[bandAndMode]
	if bandOK {
		_, result.Duplicate = callsigns[qso.TheirCall]
	}

	// find the relevant multi rules
	c.tracef("filtering %d multi rules", len(c.definition.Scoring.MultiRules))
	multiRules := c.filterScoringRules(c.definition.Scoring.MultiRules, false, c.setup.MyContinent, c.setup.MyCountry, c.setup.MyPrefix(), qso.TheirContinent, qso.TheirCountry, qso.TheirPrefix(), qso.Band, getMyProperty, getTheirProperty)
	c.tracef("found %d relevant multi rules", len(multiRules))
	for i, rule := range multiRules {
		if rule.Property == "" {
			result.Multis += rule.Value
			continue
		}

		// get the property value
		value := getTheirProperty(rule.Property)
		if value == "" {
			c.tracef("rule #%d: value is empty", i+1)
			continue
		}
		if contains(rule.Except, value) {
			c.tracef("rule #%d: value %s is excluded", i+1, value)
			continue
		}

		// apply the band rule
		bandAndMode := effectiveBandAndMode(qso.Band, qso.Mode, rule.BandRule)

		// check for duplicate values
		var duplicateValue bool
		properties, propertiesOK := c.multisPerBandAndMode[bandAndMode]
		if propertiesOK {
			c.tracef("rule #%d: multis for band and mode %v already exist: %v", i+1, bandAndMode, properties)
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

func (c *Counter) filterScoringRules(rules []ScoringRule, onlyMostRelevant bool, myContinent Continent, myCountry DXCCEntity, myPrefix string, theirContinent Continent, theirCountry DXCCEntity, theirPrefix string, band ContestBand, getMyProperty propertyProvider, getTheirProperty propertyProvider) []ScoringRule {
	matchingRules := make([]ScoringRule, 0, len(rules))
	ruleScores := make([]int, 0, len(matchingRules))
	maxRuleScores := make(map[Property]int)
	maxRuleScore := 0
	for i, rule := range rules {
		c.tracef("evaluating rule #%d:", i+1)
		ruleScore := 0

		if myContinent != "" && len(rule.MyContinent) > 0 {
			if len(rule.MyContinent) > 1 && rule.MyContinent[0] == NotContinent {
				if contains(rule.MyContinent, myContinent) {
					continue
				}
			} else if !contains(rule.MyContinent, myContinent) {
				c.tracef("not my continent %s %v", myContinent, rule.MyContinent)
				continue
			}
			ruleScore++
		}
		if myCountry != "" && len(rule.MyCountry) > 0 {
			if len(rule.MyCountry) > 1 && rule.MyCountry[0] == NotCountry {
				if contains(rule.MyCountry, myCountry) {
					continue
				}
			} else if !contains(rule.MyCountry, myCountry) {
				c.tracef("not my prefix %s %v", myPrefix, rule.MyPrefix)
				continue
			}
			ruleScore++
		}
		if myPrefix != "" && len(rule.MyPrefix) > 0 {
			if len(rule.MyPrefix) > 1 && rule.MyPrefix[0] == NotPrefix {
				if contains(rule.MyPrefix, myPrefix) {
					continue
				}
			} else if !contains(rule.MyPrefix, myPrefix) {
				c.tracef("not my country %s %v", myCountry, rule.MyCountry)
				continue
			}
			ruleScore++
		}
		if theirContinent != "" && len(rule.TheirContinent) > 0 {
			if len(rule.TheirContinent) == 1 &&
				((rule.TheirContinent[0] == SameContinent && myContinent == theirContinent) ||
					(rule.TheirContinent[0] == OtherContinent && myContinent != theirContinent)) {
				ruleScore++
			} else if len(rule.TheirContinent) > 1 && rule.TheirContinent[0] == NotContinent {
				if !contains(rule.TheirContinent, theirContinent) {
					ruleScore++
				} else {
					continue
				}
			} else if contains(rule.TheirContinent, theirContinent) {
				ruleScore++
			} else {
				c.tracef("not their continent %s %v", theirContinent, rule.TheirContinent)
				continue
			}
		}
		if theirCountry != "" && len(rule.TheirCountry) > 0 {
			if len(rule.TheirCountry) == 1 &&
				((rule.TheirCountry[0] == SameCountry && myCountry == theirCountry) ||
					(rule.TheirCountry[0] == OtherCountry && myCountry != theirCountry)) {
				ruleScore++
			} else if len(rule.TheirCountry) > 1 && rule.TheirCountry[0] == NotCountry {
				if !contains(rule.TheirCountry, theirCountry) {
					ruleScore++
				} else {
					continue
				}
			} else if contains(rule.TheirCountry, theirCountry) {
				ruleScore++
			} else {
				continue
			}
		}
		if theirPrefix != "" && len(rule.TheirPrefix) > 0 {
			if len(rule.TheirPrefix) == 1 &&
				((rule.TheirPrefix[0] == SamePrefix && myPrefix == theirPrefix) ||
					(rule.TheirPrefix[0] == OtherPrefix && myPrefix != theirPrefix)) {
				ruleScore++
			} else if len(rule.TheirPrefix) > 1 && rule.TheirPrefix[0] == NotPrefix {
				if !contains(rule.TheirPrefix, theirPrefix) {
					ruleScore++
				} else {
					continue
				}
			} else if contains(rule.TheirPrefix, theirPrefix) {
				ruleScore++
			} else {
				continue
			}
		}
		if band != "" && len(rule.Bands) > 0 {
			if contains(rule.Bands, band) {
				ruleScore++
			} else {
				c.tracef("not a valid band %s %v", band, rule.Bands)
				continue
			}
		}
		if len(rule.TheirWorkingCondition) > 0 {
			value := strings.ToLower(strings.TrimSpace(getTheirProperty(WorkingConditionProperty)))
			c.tracef("evaluating working condition %s %v", value, rule.TheirWorkingCondition)
			if len(rule.TheirWorkingCondition) == 1 && value == rule.TheirWorkingCondition[0] {
				c.tracef("working condition exact match")
				ruleScore++
			} else if len(rule.TheirWorkingCondition) > 1 && rule.TheirWorkingCondition[0] == NotPrefix {
				if !contains(rule.TheirWorkingCondition, value) {
					c.tracef("working condition match with not")
					ruleScore++
				} else {
					c.tracef("working condition NO match with not")
					continue
				}
			} else if contains(rule.TheirWorkingCondition, value) {
				c.tracef("working condition match")
				ruleScore++
			} else {
				c.tracef("working condition NO match")
				continue
			}
		}
		if rule.Property != "" {
			value := getTheirProperty(rule.Property)
			if value == "" {
				c.tracef("empty property %s", rule.Property)
				continue
			}
			ruleScore++
		}
		if len(rule.PropertyConstraints) > 0 {
			c.tracef("rule #%d has %d property constraints", i, len(rule.PropertyConstraints))
			propertyConstraintsMatched := 0

			for _, constraint := range rule.PropertyConstraints {
				myValue := getMyProperty(constraint.Name)
				theirValue := getTheirProperty(constraint.Name)
				if constraint.Matches(myValue, theirValue) {
					c.tracef("property constraint on %s matched: my %q their %q", constraint.Name, myValue, theirValue)
					propertyConstraintsMatched++
				} else {
					c.tracef("property constraint on %s NO match: my %q their %q", constraint.Name, myValue, theirValue)
				}
			}
			if propertyConstraintsMatched != len(rule.PropertyConstraints) {
				c.tracef("only %d of %d property constraints matched", propertyConstraintsMatched, len(rule.PropertyConstraints))
				continue
			}

			ruleScore++
		}

		ruleScore += rule.AdditionalWeight

		matchingRules = append(matchingRules, rule)
		ruleScores = append(ruleScores, ruleScore)
		if maxRuleScores[rule.Property] < ruleScore {
			maxRuleScores[rule.Property] = ruleScore
		}
		if maxRuleScore < ruleScore {
			maxRuleScore = ruleScore
		}
	}

	c.tracef("%d matching rules with a max score of %d", len(matchingRules), maxRuleScore)

	if maxRuleScore == 0 && len(matchingRules) > 1 {
		return []ScoringRule{}
	}

	result := make([]ScoringRule, 0, len(matchingRules))
	for i, rule := range matchingRules {
		if onlyMostRelevant {
			if ruleScores[i] == maxRuleScore {
				result = append(result, rule)
			}
		} else {
			if ruleScores[i] == maxRuleScores[rule.Property] {
				result = append(result, rule)
			}
		}
	}

	return result
}

func (c *Counter) filterExchangeFields(definitions []ExchangeDefinition, myContinent Continent, myCountry DXCCEntity, theirContinent Continent, theirCountry DXCCEntity, theirWorkingCondition string) []ExchangeField {
	theirWorkingCondition = strings.ToLower(theirWorkingCondition)
	matchingDefinitions := make([]ExchangeDefinition, 0, len(definitions))
	definitionScores := make([]int, 0, len(matchingDefinitions))
	maxDefinitionScore := 0

	for _, definition := range definitions {
		definitionScore := 0

		if myContinent != "" && len(definition.MyContinent) > 0 {
			if len(definition.MyContinent) > 1 && definition.MyContinent[0] == NotContinent {
				if contains(definition.MyContinent, myContinent) {
					continue
				}
			} else if !contains(definition.MyContinent, myContinent) {
				c.tracef("not my continent %s %v", myContinent, definition.MyContinent)
				continue
			}
			definitionScore++
		}
		if myCountry != "" && len(definition.MyCountry) > 0 {
			if len(definition.MyCountry) > 1 && definition.MyCountry[0] == NotCountry {
				if contains(definition.MyCountry, myCountry) {
					continue
				}
			} else if !contains(definition.MyCountry, myCountry) {
				c.tracef("not my country %s %v", myCountry, definition.MyCountry)
				continue
			}
			definitionScore++
		}
		if theirContinent != "" && len(definition.TheirContinent) > 0 {
			if len(definition.TheirContinent) == 1 &&
				((definition.TheirContinent[0] == SameContinent && myContinent == theirContinent) ||
					(definition.TheirContinent[0] == OtherContinent && myContinent != theirContinent)) {
				definitionScore++
			} else if len(definition.TheirContinent) > 1 && definition.TheirContinent[0] == NotContinent {
				if !contains(definition.TheirContinent, theirContinent) {
					definitionScore++
				} else {
					continue
				}
			} else if contains(definition.TheirContinent, theirContinent) {
				definitionScore++
			} else {
				c.tracef("not their continent %s %v", theirContinent, definition.TheirContinent)
				continue
			}
		}
		if theirCountry != "" && len(definition.TheirCountry) > 0 {
			if len(definition.TheirCountry) == 1 &&
				((definition.TheirCountry[0] == SameCountry && myCountry == theirCountry) ||
					(definition.TheirCountry[0] == OtherCountry && myCountry != theirCountry)) {
				definitionScore++
			} else if len(definition.TheirCountry) > 1 && definition.TheirCountry[0] == NotCountry {
				if !contains(definition.TheirCountry, theirCountry) {
					definitionScore++
				} else {
					continue
				}
			} else if contains(definition.TheirCountry, theirCountry) {
				definitionScore++
			} else {
				c.tracef("not their country %s %v", theirCountry, definition.TheirCountry)
				continue
			}
		}
		if len(definition.TheirWorkingCondition) > 0 {
			c.tracef("evaluating working condition %s %v", theirWorkingCondition, definition.TheirWorkingCondition)
			if len(definition.TheirWorkingCondition) == 1 && theirWorkingCondition == definition.TheirWorkingCondition[0] {
				c.tracef("working condition exact match")
				definitionScore++
			} else if len(definition.TheirWorkingCondition) > 1 && definition.TheirWorkingCondition[0] == NotPrefix {
				if !contains(definition.TheirWorkingCondition, theirWorkingCondition) {
					c.tracef("working condition match with not")
					definitionScore++
				} else {
					c.tracef("working condition NO match with not")
					continue
				}
			} else if contains(definition.TheirWorkingCondition, theirWorkingCondition) {
				c.tracef("working condition match")
				definitionScore++
			} else {
				c.tracef("working condition NO match")
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
