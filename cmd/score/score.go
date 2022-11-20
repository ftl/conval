package score

import (
	"fmt"

	"github.com/ftl/conval"
	"github.com/ftl/conval/app"
)

type MultisBoardRow struct {
	Property conval.Property      `yaml:"property" json:"property"`
	Multi    string               `yaml:"multi" json:"multi"`
	Bands    []conval.ContestBand `yaml:"bands,flow" json:"bands"`
}

type Result struct {
	MultiProperties []conval.Property    `yaml:"multi_properties,flow" json:"multi_properties"`
	Bands           []conval.ContestBand `yaml:"bands,flow" json:"bands"`

	MultisBoard []MultisBoardRow `yaml:"multis_board" json:"multis_board"`
	QSOs        int              `yaml:"qsos" json:"qsos"`
	Points      int              `yaml:"points" json:"points"`
	Multis      int              `yaml:"multis" json:"multis"`
	Total       int              `yaml:"total" json:"total"`
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

	counter := conval.NewCounter(*setupForFile, definitionForFile.Exchange, definitionForFile.Scoring)
	qsos := logfile.QSOs(counter.EffectiveExchangeFields)
	for _, qso := range qsos {
		counter.Add(qso)
	}

	multiProperties := counter.MultiProperties()
	bands := counter.UsedBands()
	multisBoard := make([]MultisBoardRow, 0, len(multiProperties)*len(bands))
	for _, property := range multiProperties {
		for _, multi := range counter.MultisPerProperty(property) {
			bandsPerMulti := counter.BandsPerMulti(property, multi)
			multisBoard = append(multisBoard, MultisBoardRow{property, multi, bandsPerMulti})
		}
	}
	totalScore := counter.TotalScore()
	return Result{
		MultiProperties: multiProperties,
		Bands:           bands,
		MultisBoard:     multisBoard,
		QSOs:            totalScore.QSOs,
		Points:          totalScore.Points,
		Multis:          totalScore.Multis,
		Total:           totalScore.Total(),
	}, nil
}
