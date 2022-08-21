package conval

type Counter struct {
	setup Setup
	rules Scoring

	scorePerBand map[ContestBand]BandScore
	overallScore BandScore
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
	}
}

func (c *Counter) Add(qso QSO) QSOScore {
	return QSOScore{} // TODO implement
}

func (c Counter) Probe(qso QSO) QSOScore {
	return QSOScore{} // TODO implement
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
