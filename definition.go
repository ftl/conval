package conval

import (
	"fmt"
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

type Definition struct {
	Name          string                `yaml:"name"`
	Identifier    string                `yaml:"identifier"`
	OfficialRules string                `yaml:"official_rules"`
	Durations     []ConstrainedDuration `yaml:"durations"`
	Breaks        []ConstrainedDuration `yaml:"breaks"`
	Categories    []Category            `yaml:"categories"`
	Overlays      []Overlay             `yaml:"overlays"`
	Modes         []Mode                `yaml:"modes"`
	Bands         []ContestBand         `yaml:"bands"`
}

type ConstrainedDuration struct {
	Constraint `yaml:",inline"`
	Duration   time.Duration `yaml:"duration"`
}

type Category struct {
	Name     string       `yaml:"name"`
	Operator OperatorMode `yaml:"operator"`
	TX       TXMode       `yaml:"tx"`
	Power    PowerMode    `yaml:"power"`
	Bands    BandMode     `yaml:"bands"`
	Modes    []Mode       `yaml:"modes"`
	Assisted bool         `yaml:"assisted"`
}

func LoadYAML(r io.Reader) (*Definition, error) {
	decoder := yaml.NewDecoder(r)

	var result Definition
	err := decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	result.Sanitize()

	return &result, nil
}

func (d *Definition) Sanitize() {
	// TODO implement
	// - expand the all keyword for modes and bands
	// - make all enum values lower case to match the defined constants
}

func (d Definition) Validate() error {
	// TODO implement
	return fmt.Errorf("not yet implemented")
}
