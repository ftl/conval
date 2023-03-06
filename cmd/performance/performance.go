package performance

import (
	"fmt"
	"time"

	"github.com/ftl/conval"
	"github.com/ftl/conval/app"
)

type DataPoint struct {
	Offset      time.Duration `yaml:"offset" json:"offset"`
	Bands       []string      `yaml:"bands,flow" json:"bands"`
	MultiValues []string      `yaml:"multi_values" json:"multi_values"`
	QSOs        int           `yaml:"qsos" json:"qsos"`
	Points      int           `yaml:"points" json:"points"`
	Multis      int           `yaml:"multis" json:"multis"`
	Total       int           `yaml:"total" json:"total"`
}

type Result struct {
	Resolution time.Duration
	DataPoints []DataPoint
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

	counter := conval.NewCounter(*definitionForFile, *setupForFile)
	qsos := logfile.QSOs(definition, counter.EffectiveExchangeFields)
	for _, qso := range qsos {
		counter.Add(qso)
	}

	scoreBins := counter.EvaluateAll(startTime, resolution)
	result := Result{
		Resolution: resolution,
		DataPoints: make([]DataPoint, len(scoreBins)),
	}
	for i, bin := range scoreBins {
		bands := make([]string, 0, len(bin.Bands))
		for band, used := range bin.Bands {
			if used {
				bands = append(bands, string(band))
			}
		}
		// TODO sort bands

		multiValues := make([]string, 0, 100)
		for _, multisPerProperty := range bin.MultiValues {
			for multi := range multisPerProperty {
				multiValues = append(multiValues, multi)
			}
		}

		result.DataPoints[i] = DataPoint{
			Offset:      time.Duration(i+1) * resolution,
			Bands:       bands,
			MultiValues: multiValues,
			QSOs:        bin.QSOs,
			Points:      bin.Points,
			Multis:      bin.Multis,
			Total:       bin.Total(),
		}
	}

	return result, nil
}
