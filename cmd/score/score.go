package score

import (
	"fmt"

	"github.com/ftl/conval"
)

type Logfile interface {
	Identifier() conval.ContestIdentifier
	Setup() *conval.Setup
	QSOs(exchangeFields func(conval.Continent, conval.DXCCEntity) []conval.ExchangeField) []conval.QSO
}

type MultisBoardRow struct {
	Property conval.Property      `yaml:"property" json:"property"`
	Multi    string               `yaml:"multi" json:"multi"`
	Bands    []conval.ContestBand `yaml:"bands" json:"bands"`
}

type Result struct {
	MultiProperties []conval.Property    `yaml:"multi_properties" json:"multi_properties"`
	Bands           []conval.ContestBand `yaml:"bands" json:"bands"`

	MultisBoard []MultisBoardRow `yaml:"multis_board" json:"multis_board"`
	QSOs        int              `yaml:"qsos" json:"qsos"`
	Points      int              `yaml:"points" json:"points"`
	Multis      int              `yaml:"multis" json:"multis"`
	Total       int              `yaml:"total" json:"total"`
}

func PrepareDefinition(filename string, cabrilloName string) (*conval.Definition, error) {
	if filename != "" {
		return conval.LoadDefinitionFromFile(filename)
	}
	if cabrilloName != "" {
		return conval.IncludedDefinition(cabrilloName)
	}
	return nil, nil
}

func PrepareSetup(filename string, prefixes conval.PrefixDatabase) (*conval.Setup, error) {
	// TODO add parameters for things that can be overridden with CLI flags
	var result *conval.Setup
	var err error
	if filename != "" {
		result, err = conval.LoadSetupFromFile(filename)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	myContinent, myCountry, found := prefixes.Find(result.MyCall.String())
	if found && result.MyContinent == "" {
		result.MyContinent = myContinent
	}
	if found && result.MyCountry == "" {
		result.MyCountry = myCountry
	}

	return result, nil
}

func Evaluate(logfile Logfile, definition *conval.Definition, setup *conval.Setup) (Result, error) {
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
