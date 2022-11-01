package conval

import (
	"fmt"
)

func ValidateExamples(definition *Definition) error {
	for i, example := range definition.Examples {
		// log.Printf("example #%d", i+1)
		err := validateExample(definition, example)
		if err != nil {
			return fmt.Errorf("example #%d is invalid: %w", i+1, err)
		}
	}
	return nil
}

func validateExample(definition *Definition, example Example) error {
	counter := NewCounter(example.Setup.ToSetup(), definition.Exchange, definition.Scoring)
	for i, qso := range example.QSOs {
		exchangeFields := counter.EffectiveExchangeFields(qso.TheirContinent, qso.TheirCountry)
		// log.Printf("QSO #%d with exchange fields %v", i+1, exchangeFields)
		qsoScore := counter.Add(qso.ToQSO(exchangeFields))
		if !(qso.Score.Equal(qsoScore)) {
			return fmt.Errorf("the score of QSO #%d is wrong, expected %d points * %d multis, duplicate should be %t, but got %d points * %d multis, duplicate is %t", i+1, qso.Score.Points, qso.Score.Multis, qso.Score.Duplicate, qsoScore.Points, qsoScore.Multis, qsoScore.Duplicate)
		}
	}

	totalScore := counter.TotalScore()
	if example.Score != totalScore {
		return fmt.Errorf("the total score is wrong, expected %d points * %d multis, but got %d points * %d multis", example.Score.Points, example.Score.Multis, totalScore.Points, totalScore.Multis)
	}

	return nil
}
