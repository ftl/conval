package app

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ftl/conval"
)

type Logfile interface {
	Identifier() conval.ContestIdentifier
	Setup() *conval.Setup
	QSOs(exchangeFields func(conval.Continent, conval.DXCCEntity) []conval.ExchangeField) []conval.QSO
}

type OutputFormat string

const (
	TextOutput OutputFormat = "text"
	YamlOutput OutputFormat = "yaml"
	JsonOutput OutputFormat = "json"
	CsvOutput  OutputFormat = "csv"
)

func ParseOutputFormat(s string) OutputFormat {
	return OutputFormat(strings.ToLower(strings.TrimSpace(s)))
}

func PrepareDefinition(name string) (*conval.Definition, error) {
	if name == "" {
		return nil, nil
	}

	result, err := conval.LoadDefinitionFromFile(name)
	if err == nil {
		return result, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	result, err = conval.IncludedDefinition(name)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("the contest definition %s does not exist", name)
	}
	return result, err
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
