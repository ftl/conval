package conval

import (
	"fmt"
	"log"
)

func ValidateExamples(definition *Definition, prefixes PrefixDatabase) error {
	for i, example := range definition.Examples {
		err := validateExample(definition, example, prefixes, false)
		if err != nil {
			return fmt.Errorf("example #%d is invalid: %w", i+1, err)
		}
	}
	return nil
}

func ValidateExamplesTrace(definition *Definition, prefixes PrefixDatabase) error {
	for i, example := range definition.Examples {
		log.Printf("example #%d", i+1)
		err := validateExample(definition, example, prefixes, true)
		if err != nil {
			return fmt.Errorf("example #%d is invalid: %w", i+1, err)
		}
	}
	return nil
}

func validateExample(definition *Definition, example Example, prefixes PrefixDatabase, trace bool) error {
	counter := NewCounter(*definition, example.Setup.ToSetup())
	counter.SetTrace(trace)
	for i, qso := range example.QSOs {
		exchangeFields := counter.EffectiveExchangeFields(qso.TheirContinent, qso.TheirCountry)
		qsoScore := counter.Add(qso.ToQSO(exchangeFields, example.Setup.MyExchange, prefixes, definition))
		if trace {
			log.Printf("QSO #%d with exchange fields %v: %+v", i+1, exchangeFields, qsoScore)
		}
		if !(qso.Score.Equal(qsoScore)) {
			return fmt.Errorf("the score of QSO #%d is wrong, expected %d points * %d multis, duplicate should be %t, but got %d points * %d multis, duplicate is %t", i+1, qso.Score.Points, qso.Score.Multis, qso.Score.Duplicate, qsoScore.Points, qsoScore.Multis, qsoScore.Duplicate)
		}
	}

	totalScore := counter.TotalScore()
	if trace {
		log.Printf("%+v = %d", totalScore, counter.Total(totalScore))
	}
	if !equalScore(counter, example.Score, totalScore) {
		return fmt.Errorf("the total score is wrong, expected %d qsos with %d points * %d multis, but got %d qsos with %d points * %d multis", example.Score.QSOs, example.Score.Points, example.Score.Multis, totalScore.QSOs, totalScore.Points, totalScore.Multis)
	}

	return nil
}

func equalScore(counter *Counter, expected ScoreExample, actual BandScore) bool {
	return expected.QSOs == actual.QSOs &&
		expected.Points == actual.Points &&
		expected.Multis == actual.Multis &&
		expected.Total == counter.Total(actual)
}
