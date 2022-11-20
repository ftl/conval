package performance

import (
	"fmt"
	"time"

	"github.com/ftl/conval"
	"github.com/ftl/conval/app"
)

type Result struct {
}

func Evaluate(logfile app.Logfile, definition *conval.Definition, setup *conval.Setup, startTime time.Time, resolution time.Duration) (Result, error) {
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

	counter := conval.NewCounter(*setupForFile, definitionForFile.Exchange, definitionForFile.Scoring)
	qsos := logfile.QSOs(counter.EffectiveExchangeFields)
	for _, qso := range qsos {
		counter.Add(qso)
	}

	return Result{}, fmt.Errorf("not yet implemented")
}
