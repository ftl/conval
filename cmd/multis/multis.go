package multis

import (
	"fmt"

	"github.com/ftl/conval"
	"github.com/ftl/conval/app"
)

type Row struct {
	Property conval.Property      `yaml:"property" json:"property"`
	Multi    string               `yaml:"multi" json:"multi"`
	Bands    []conval.ContestBand `yaml:"bands,flow" json:"bands"`
}

type Result struct {
	MultiProperties []conval.Property    `yaml:"multi_properties,flow" json:"multi_properties"`
	Bands           []conval.ContestBand `yaml:"bands,flow" json:"bands"`
	Rows            []Row                `yaml:"rows" json:"rows"`
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
	qsos := logfile.QSOs(definition, counter.EffectiveExchangeFields)
	for _, qso := range qsos {
		counter.Add(qso)
	}

	multiProperties := counter.MultiProperties()
	bands := counter.UsedBands()
	rows := make([]Row, 0, len(multiProperties)*len(bands))
	for _, property := range multiProperties {
		for _, multi := range counter.MultisPerProperty(property) {
			bandsPerMulti := counter.BandsPerMulti(property, multi)
			rows = append(rows, Row{property, multi, bandsPerMulti})
		}
	}
	result := Result{
		MultiProperties: multiProperties,
		Bands:           bands,
		Rows:            rows,
	}

	return result, nil
}
