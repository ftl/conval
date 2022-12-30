package conval

import (
	"io"
	"os"
	"time"

	"github.com/ftl/hamradio/callsign"
	"github.com/ftl/hamradio/locator"
	"gopkg.in/yaml.v3"
)

// A ContestIdentifier aka. Cabrillo name is a unique identifier for a contest.
// See https://www.contestcalendar.com/cabnames.php
type ContestIdentifier string

// A Definition of a contest.
type Definition struct {
	Name                string                `yaml:"name"`
	Identifier          ContestIdentifier     `yaml:"identifier"`
	OfficialRules       string                `yaml:"official_rules"`
	UploadURL           string                `yaml:"upload_url"`
	UploadFormat        string                `yaml:"upload_format"`
	Duration            time.Duration         `yaml:"duration,omitempty"`
	DurationConstraints []ConstrainedDuration `yaml:"duration-constraints,omitempty"`
	Breaks              []ConstrainedDuration `yaml:"breaks,omitempty"`
	Categories          []Category            `yaml:"categories,omitempty"`
	Overlays            []Overlay             `yaml:"overlays,omitempty"`
	Modes               []Mode                `yaml:"modes,omitempty"`
	Bands               []ContestBand         `yaml:"bands,omitempty"`
	BandChangeRules     []BandChangeRule      `yaml:"band_change_rules,omitempty"`
	Exchange            []ExchangeDefinition  `yaml:"exchange"`
	Scoring             Scoring               `yaml:"scoring"`
	Examples            []Example             `yaml:"examples,omitempty"`
}

func (d Definition) ExchangeFields() []ExchangeField {
	result := make([]ExchangeField, 0, 3)
	usedProperties := make([]map[Property]bool, 0, 3)
	for _, definition := range d.Exchange {
		for i, field := range definition.Fields {
			if i >= len(result) {
				result = append(result, field)
				usedProperties = append(usedProperties, make(map[Property]bool))
				for _, property := range field {
					usedProperties[i][property] = true
				}
				continue
			}

			for _, property := range field {
				if usedProperties[i][property] {
					continue
				}
				result[i] = append(result[i], property)
				usedProperties[i][property] = true
			}
		}
	}
	return result
}

type ConstrainedDuration struct {
	Constraint `yaml:",inline"`
	Duration   time.Duration `yaml:"duration"`
}

type BandChangeRule struct {
	Constraint          `yaml:",inline"`
	GracePeriod         time.Duration `yaml:"grace_period"`
	MultiplierException bool          `yaml:"multiplier_exception"`
}

type Category struct {
	Name     string       `yaml:"name"`
	Operator OperatorMode `yaml:"operator,omitempty"`
	TX       TXMode       `yaml:"tx,omitempty"`
	Power    PowerMode    `yaml:"power,omitempty"`
	Bands    BandMode     `yaml:"bands,omitempty"`
	Modes    []Mode       `yaml:"modes,omitempty"`
	Assisted bool         `yaml:"assisted,omitempty"`
	Overlay  Overlay      `yaml:"overlay,omitempty"`
}

type ExchangeDefinition struct {
	MyContinent           []Continent     `yaml:"my_continent,omitempty"`
	MyCountry             []DXCCEntity    `yaml:"my_country,omitempty"`
	TheirContinent        []Continent     `yaml:"their_continent,omitempty"`
	TheirCountry          []DXCCEntity    `yaml:"their_country,omitempty"`
	TheirWorkingCondition string          `yaml:"their_working_condition,omitempty"`
	AdditionalWeight      int             `yaml:"additional_weight,omitempty"`
	Fields                []ExchangeField `yaml:"fields,omitempty"`
}

type ExchangeField []Property

func (f ExchangeField) Strings() []string {
	result := make([]string, len(f))
	for i, p := range f {
		result[i] = string(p)
	}
	return result
}

func (f ExchangeField) Contains(property Property) bool {
	for _, p := range f {
		if p == property {
			return true
		}
	}
	return false
}

type Scoring struct {
	QSORules    []ScoringRule `yaml:"qsos"`
	QSOBandRule BandRule      `yaml:"qso_band_rule"`
	MultiRules  []ScoringRule `yaml:"multis"`
}

type ScoringRule struct {
	MyContinent           []Continent   `yaml:"my_continent,omitempty"`
	MyCountry             []DXCCEntity  `yaml:"my_country,omitempty"`
	TheirContinent        []Continent   `yaml:"their_continent,omitempty"`
	TheirCountry          []DXCCEntity  `yaml:"their_country,omitempty"`
	TheirWorkingCondition string        `yaml:"their_working_condition,omitempty"`
	Bands                 []ContestBand `yaml:"bands,omitempty"`
	Property              Property      `yaml:"property,omitempty"`
	BandRule              BandRule      `yaml:"band_rule,omitempty"`
	AdditionalWeight      int           `yaml:"additional_weight,omitempty"`
	Value                 int           `yaml:"value,omitempty"`
}

type Example struct {
	Setup SetupExample `yaml:"setup"`
	QSOs  []QSOExample `yaml:"qsos"`
	Score BandScore    `yaml:"score"`
}

type SetupExample struct {
	MyCall      string     `yaml:"my_call,omitempty"`
	MyContinent Continent  `yaml:"my_continent,omitempty"`
	MyCountry   DXCCEntity `yaml:"my_country,omitempty"`

	GridLocator string   `yaml:"grid_locator,omitempty"`
	Operators   []string `yaml:"operators,omitempty"`

	OperatorMode OperatorMode  `yaml:"operator_mode,omitempty"`
	Overlay      Overlay       `yaml:"overlay,omitempty"`
	Power        PowerMode     `yaml:"power,omitempty"`
	Bands        []ContestBand `yaml:"bands,omitempty"`
	Modes        []Mode        `yaml:"modes,omitempty"`

	MyExchange QSOExchange `yaml:"my_exchange,omitempty"`
}

func (s SetupExample) ToSetup() Setup {
	myCall, err := callsign.Parse(s.MyCall)
	if err != nil {
		myCall = callsign.Callsign{}
	}
	gridLocator, err := locator.Parse(s.GridLocator)
	if err != nil {
		gridLocator = locator.Locator{}
	}
	operators := make([]callsign.Callsign, 0, len(s.Operators))
	for _, operator := range s.Operators {
		operatorCall, err := callsign.Parse(operator)
		if err == nil {
			operators = append(operators, operatorCall)
		}
	}

	return Setup{
		MyCall:       myCall,
		MyContinent:  s.MyContinent,
		MyCountry:    s.MyCountry,
		GridLocator:  gridLocator,
		Operators:    operators,
		OperatorMode: s.OperatorMode,
		Overlay:      s.Overlay,
		Power:        s.Power,
		Bands:        s.Bands,
		Modes:        s.Modes,
		MyExchange:   s.MyExchange,
	}
}

type QSOExample struct {
	TheirCall      string     `yaml:"their_call,omitempty"`
	TheirContinent Continent  `yaml:"their_continent,omitempty"`
	TheirCountry   DXCCEntity `yaml:"their_country,omitempty"`

	Timestamp time.Time   `yaml:"time,omitempty"`
	Band      ContestBand `yaml:"band,omitempty"`
	Mode      Mode        `yaml:"mode,omitempty"`

	MyExchange    QSOExchange `yaml:"my_exchange,omitempty"`
	TheirExchange []string    `yaml:"their_exchange,omitempty"`

	Score QSOScore `yaml:",inline"`
}

func (q QSOExample) ToQSO(fields []ExchangeField, prefixes PrefixDatabase) QSO {
	return QSO{
		TheirCall:      callsign.MustParse(q.TheirCall),
		TheirContinent: q.TheirContinent,
		TheirCountry:   q.TheirCountry,
		Timestamp:      q.Timestamp,
		Band:           q.Band,
		Mode:           q.Mode,
		MyExchange:     q.MyExchange,
		TheirExchange:  ParseExchange(fields, q.TheirExchange, prefixes),
	}
}

func LoadDefinitionFromFile(filename string) (*Definition, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return LoadDefinitionYAML(file)
}

func LoadDefinitionYAML(r io.Reader) (*Definition, error) {
	decoder := yaml.NewDecoder(r)

	var result Definition
	err := decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func SaveDefinitionYAML(w io.Writer, definition *Definition, withExamples bool) error {
	if definition == nil {
		return nil
	}

	encoder := yaml.NewEncoder(w)

	if withExamples {
		return encoder.Encode(definition)
	}

	definitionWithoutExamples := *definition
	definitionWithoutExamples.Examples = nil
	return encoder.Encode(definitionWithoutExamples)
}

func LoadSetupFromFile(filename string) (*Setup, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return LoadSetupYAML(file)
}

func LoadSetupYAML(r io.Reader) (*Setup, error) {
	decoder := yaml.NewDecoder(r)

	var setup SetupExample
	err := decoder.Decode(&setup)
	if err != nil {
		return nil, err
	}

	result := setup.ToSetup()
	return &result, nil
}
