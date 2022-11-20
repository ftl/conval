package score

import (
	"fmt"

	"github.com/ftl/conval"
	"github.com/ftl/conval/app"
)

type Result struct {
	QSOs   int `yaml:"qsos" json:"qsos"`
	Points int `yaml:"points" json:"points"`
	Multis int `yaml:"multis" json:"multis"`
	Total  int `yaml:"total" json:"total"`
}

func Evaluate(logfile app.Logfile, definition *conval.Definition, setup *conval.Setup) (Result, error) {
	var err error

	definitionForFile := definition
	if definitionForFile == nil {
		definitionForFile, err = conval.IncludedDefinition(string(logfile.Identifier()))
		if err != nil {
			return Result{}, err
		}
	}
	if definitionForFile == nil {
		return Result{}, fmt.Errorf("no contest definition found")
	}

	setupForFile := setup
	if setupForFile == nil {
		setupForFile = logfile.Setup()
	}
	if setupForFile == nil {
		return Result{}, fmt.Errorf("no setup defined")
	}

	counter := conval.NewCounter(*definitionForFile, *setupForFile)
	qsos := logfile.QSOs(counter.EffectiveExchangeFields)
	for _, qso := range qsos {
		counter.Add(qso)
	}

	totalScore := counter.TotalScore()
	return Result{
		QSOs:   totalScore.QSOs,
		Points: totalScore.Points,
		Multis: totalScore.Multis,
		Total:  totalScore.Total(),
	}, nil
}
