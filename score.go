package conval

import "github.com/ftl/hamradio/callsign"

type Counter struct {
	setup Setup
	rules Scoring

	callsignsPerBand map[ContestBand]map[callsign.Callsign]bool
	scorePerBand     map[ContestBand]BandScore
}

type QSOScore struct {
	Points    int
	Multis    int
	Duplicate bool
}

type BandScore struct {
	Points int
	Multis int
}

func NewCounter(setup Setup, rules Scoring) *Counter {
	return &Counter{
		setup: setup,
		rules: rules,

		callsignsPerBand: make(map[ContestBand]map[callsign.Callsign]bool),
		scorePerBand:     make(map[ContestBand]BandScore),
	}
}

func (c Counter) TotalScore() BandScore {
	return c.scorePerBand[BandAll]
}

func (c *Counter) Add(qso QSO) QSOScore {
	result := c.Probe(qso)

	// apply the QSO band rule
	var band ContestBand
	switch c.rules.QSOBandRule {
	case Once:
		band = BandAll
	case OncePerBand:
		band = qso.Band
	}

	// update the callsign registry
	var callsigns map[callsign.Callsign]bool
	var bandOK bool
	callsigns, bandOK = c.callsignsPerBand[band]
	if !bandOK {
		callsigns = make(map[callsign.Callsign]bool)
	}
	callsigns[qso.TheirCall] = true
	c.callsignsPerBand[band] = callsigns

	// update the scores
	if !result.Duplicate {
		totalScore := c.scorePerBand[BandAll]
		totalScore.Points += result.Points
		// TODO: add multis
		c.scorePerBand[BandAll] = totalScore

		scorePerBand := c.scorePerBand[qso.Band]
		scorePerBand.Points += result.Points
		// TODO: add multis
		c.scorePerBand[qso.Band] = scorePerBand
	}

	// TODO: update other internal data structures for duplicate checks etc.

	return result
}

func (c Counter) Probe(qso QSO) QSOScore {
	result := QSOScore{}

	qsoRules := filterScoringRules(c.rules.QSORules, true, c.setup.MyContinent, c.setup.MyCountry, qso.TheirContinent, qso.TheirCountry, qso.Band, qso.TheirExchange)
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
	var band ContestBand
	switch c.rules.QSOBandRule {
	case Once:
		band = BandAll
	case OncePerBand:
		band = qso.Band
	}
	// check the callsign registry for duplicate
	var callsigns map[callsign.Callsign]bool
	var bandOK bool
	callsigns, bandOK = c.callsignsPerBand[band]
	if bandOK {
		_, result.Duplicate = callsigns[qso.TheirCall]
	}

	multiRules := filterScoringRules(c.rules.MultiRules, true, c.setup.MyContinent, c.setup.MyCountry, qso.TheirContinent, qso.TheirCountry, qso.Band, qso.TheirExchange)
	for _, rule := range multiRules {
		// TODO:
		// - check for uniqueness
		// - apply the band rule
		result.Multis += rule.Value
	}

	return result
}

func filterScoringRules(rules []ScoringRule, onlyMostRelevant bool, myContinent Continent, myCountry DXCCEntity, theirContinent Continent, theirCountry DXCCEntity, band ContestBand, exchange QSOExchange) []ScoringRule {
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
			if exchange == nil {
				continue
			}
			_, ok := exchange[rule.Property]
			if !ok {
				continue
			}
			ruleScore++
		}

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
