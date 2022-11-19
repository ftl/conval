package score

import (
	"fmt"

	"github.com/ftl/conval"
)

type Logfile interface {
	Identifier() conval.ContestIdentifier
	Setup() *conval.Setup
	QSOs() []conval.QSO
}

type Result struct {
	QSOs   int `yaml:"qsos" json:"qsos"`
	Points int `yaml:"points" json:"points"`
	Multis int `yaml:"multis" json:"multis"`
	Total  int `yaml:"total" json:"total"`
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

func Evaluate(logfile Logfile, setup *conval.Setup, definition *conval.Definition) (Result, error) {
	return Result{}, fmt.Errorf("not yet implemented")
}
