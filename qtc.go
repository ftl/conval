package conval

type ScoredQTC struct {
	QTC
	QTCScore
}

type QTCScore struct {
	Value int `yaml:"value" json:"value"`
}

func (c *Counter) AddQTC(qtc QTC) QTCScore {
	result := c.ProbeQTC(qtc)
	c.qtcs = append(c.qtcs, ScoredQTC{QTC: qtc, QTCScore: result})

	// update the scores
	totalScore := c.scorePerBand[BandAll]
	totalScore.QTCs += result.Value
	c.scorePerBand[BandAll] = totalScore

	scorePerBand := c.scorePerBand[qtc.Band]
	scorePerBand.QTCs += result.Value
	c.scorePerBand[qtc.Band] = scorePerBand

	return result
}

func (c *Counter) ProbeQTC(qtc QTC) QTCScore {
	tracef("probing QTC %+v", qtc)

	result := QTCScore{}

	getMyProperty := func(property Property) string {
		getter, getterOK := c.definition.MyQTCPropertyGetter(property)
		if !getterOK {
			tracef("no QTC property getter for my %s", property)
			return ""
		}
		return getter.GetQTCProperty(qtc, c.setup, c.prefixes)
	}

	getTheirProperty := func(property Property) string {
		getter, getterOK := c.definition.QTCPropertyGetter(property)
		if !getterOK {
			tracef("no QTC property getter for %s", property)
			return ""
		}
		return getter.GetQTCProperty(qtc, c.setup, c.prefixes)
	}

	// find the relevant QTC rules
	tracef("filtering %d QTC scoring rules", len(c.definition.Scoring.QTCRules))
	qtcRules := c.filterScoringRules(c.definition.Scoring.QSORules, true, c.setup.MyContinent, c.setup.MyCountry, c.setup.MyPrefix(), qtc.TheirContinent, qtc.TheirCountry, qtc.TheirPrefix(), qtc.Band, getMyProperty, getTheirProperty)
	tracef("found %d relevant QTC rules: %+v", len(qtcRules), qtcRules)
	if len(qtcRules) == 1 {
		result.Value = qtcRules[0].Value * qtc.Count
	} else if len(qtcRules) > 1 {
		value := qtcRules[0].Value * qtc.Count
		allEqual := true
		for _, rule := range qtcRules {
			if value != (rule.Value * qtc.Count) {
				allEqual = false
				break
			}
		}
		if allEqual {
			result.Value = value
		} else {
			tracef("inconsistent QTC rules for QTC %s with %s (%s, %s): %+v", qtc.Header, qtc.TheirCall, qtc.TheirContinent, qtc.TheirCountry, qtcRules)
		}
	}

	return result
}
