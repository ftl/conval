package conval

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type Definition struct {
	Name          string `yaml:"name"`
	Identifier    string `yaml:"identifier"`
	OfficialRules string `yaml:"official_rules"`
}

func LoadYAML(r io.Reader) (*Definition, error) {
	decoder := yaml.NewDecoder(r)

	var result Definition
	err := decoder.Decode(&result)
	if err != nil {
		return nil, err
	}
	// TODO result.Validate
	return &result, nil
}

func (d Definition) Validate() error {
	// TODO implement
	return fmt.Errorf("not yet implemented")
}
