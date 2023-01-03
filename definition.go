package conval

import (
	"io"
	"os"
	"strconv"
	"strings"
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
	fieldCount := 0
	for _, definition := range d.Exchange {
		definitionFieldCount := len(definition.Fields)
		if fieldCount < definitionFieldCount {
			fieldCount = definitionFieldCount
		}
	}
	result := make([]ExchangeField, fieldCount)
	usedProperties := make([]map[Property]bool, fieldCount)
	appendProperty := func(field ExchangeField, usedProperties map[Property]bool, property Property) (ExchangeField, map[Property]bool) {
		if usedProperties[property] {
			return field, usedProperties
		}
		usedProperties[property] = true
		field = append(field, property)
		return field, usedProperties
	}
	for _, definition := range d.Exchange {
		for i := range result {
			if usedProperties[i] == nil {
				usedProperties[i] = make(map[Property]bool)
			}
			if i >= len(definition.Fields) {
				result[i], usedProperties[i] = appendProperty(result[i], usedProperties[i], EmptyProperty)
				continue
			}

			field := definition.Fields[i]
			for _, property := range field {
				result[i], usedProperties[i] = appendProperty(result[i], usedProperties[i], property)
			}
		}
	}
	return result
}

type ConstrainedDuration struct {
	Constraint `yaml:",inline"`
	Duration   time.Duration          `yaml:"duration"`
	Mode       DurationConstraintMode `yaml:"constraint_mode,omitempty"`
}

type DurationConstraintMode string

const (
	// TotalTime counts from the timestamp of the first QSO until the timestamp of the last QSO without considering breaks.
	TotalTime DurationConstraintMode = "total_time"
	// ActiveTime counts from the timestamp of the first QSO until the timestamp of the last QSO, breaks that last at least one hour are subtracted.
	ActiveTime DurationConstraintMode = "active_time"
)

type BandChangeRule struct {
	Constraint          `yaml:",inline"`
	GracePeriod         time.Duration `yaml:"grace_period"`
	MultiplierException bool          `yaml:"multiplier_exception"`
}

type Category struct {
	Name      string       `yaml:"name"`
	Operator  OperatorMode `yaml:"operator,omitempty"`
	TX        TXMode       `yaml:"tx,omitempty"`
	Power     PowerMode    `yaml:"power,omitempty"`
	Bands     BandMode     `yaml:"bands,omitempty"`
	Modes     []Mode       `yaml:"modes,omitempty"`
	Assisted  bool         `yaml:"assisted,omitempty"`
	Overlay   Overlay      `yaml:"overlay,omitempty"`
	ScoreMode ScoreMode    `yaml:"score_mode,omitempty"`
}

type ScoreMode string

const (
	// StrictScore allows only the number of bands defined in the category (single, <any number>, all). If more bands were worked, the claimed score is zero.
	StrictScore ScoreMode = "strict"
	// BestScore counts only the best n bands, according to the number of bands defined in the category (single, <any number>, all).
	BestScore ScoreMode = "best"
)

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
	QSORules       []ScoringRule  `yaml:"qsos,omitempty"`
	QSOBandRule    BandRule       `yaml:"qso_band_rule,omitempty"`
	MultiRules     []ScoringRule  `yaml:"multis,omitempty"`
	MultiOperation MultiOperation `yaml:"multi_operation,omitempty"`
}

type ScoringRule struct {
	MyContinent           []Continent          `yaml:"my_continent,omitempty"`
	MyCountry             []DXCCEntity         `yaml:"my_country,omitempty"`
	TheirContinent        []Continent          `yaml:"their_continent,omitempty"`
	TheirCountry          []DXCCEntity         `yaml:"their_country,omitempty"`
	TheirWorkingCondition string               `yaml:"their_working_condition,omitempty"`
	Bands                 []ContestBand        `yaml:"bands,omitempty"`
	Property              Property             `yaml:"property,omitempty"` // only useful for multis
	PropertyConstraints   []PropertyConstraint `yaml:"property_constraints,omitempty"`
	BandRule              BandRule             `yaml:"band_rule,omitempty"`
	AdditionalWeight      int                  `yaml:"additional_weight,omitempty"`
	Value                 int                  `yaml:"value,omitempty"`
}

type MultiOperation string

const (
	DefaultMultiOperation MultiOperation = ""
	MultiplyMultis        MultiOperation = "multiply"
	AddMultis             MultiOperation = "add"
)

type PropertyConstraint struct {
	Name       Property `yaml:"name"`
	Min        string   `yaml:"min,omitempty"`
	Max        string   `yaml:"max,omitempty"`
	MyValue    string   `yaml:"my_value,omitempty"`
	TheirValue string   `yaml:"their_value,omitempty"`
}

func (c PropertyConstraint) Matches(myValue string, theirValue string) bool {
	myValue = sanitizePropertyValue(myValue)
	theirValue = sanitizePropertyValue(theirValue)
	result := true
	if c.Min != "" || c.Max != "" {
		result = result && c.matchesMinMax(theirValue)
	}
	if c.MyValue != "" {
		result = result && myValue == c.MyValue
	}
	if c.TheirValue != "" {
		result = result && theirValue == c.TheirValue
	}
	return result
}

func (c PropertyConstraint) matchesMinMax(value string) bool {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	if c.Min != "" && c.Max != "" {
		min, err := strconv.Atoi(c.Min)
		if err != nil {
			return false
		}
		max, err := strconv.Atoi(c.Max)
		if err != nil {
			return false
		}
		if intValue < min || intValue > max {
			return false
		}
	} else if c.Min != "" {
		min, err := strconv.Atoi(c.Min)
		if err != nil {
			return false
		}
		if intValue < min {
			return false
		}
	} else if c.Max != "" {
		max, err := strconv.Atoi(c.Max)
		if err != nil {
			return false
		}
		if intValue > max {
			return false
		}
	}
	return true
}

func sanitizePropertyValue(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

type Example struct {
	Setup SetupExample `yaml:"setup"`
	QSOs  []QSOExample `yaml:"qsos"`
	Score ScoreExample `yaml:"score"`
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

	TheirExchange []string `yaml:"their_exchange,omitempty"`

	Score QSOScore `yaml:",inline"`
}

func (q QSOExample) ToQSO(fields []ExchangeField, myExchange QSOExchange, prefixes PrefixDatabase) QSO {
	return QSO{
		TheirCall:      callsign.MustParse(q.TheirCall),
		TheirContinent: q.TheirContinent,
		TheirCountry:   q.TheirCountry,
		Timestamp:      q.Timestamp,
		Band:           q.Band,
		Mode:           q.Mode,
		MyExchange:     myExchange,
		TheirExchange:  ParseExchange(fields, q.TheirExchange, prefixes),
	}
}

type ScoreExample struct {
	QSOs   int `yaml:"qsos,omitempty"`
	Points int `yaml:"points,omitempty"`
	Multis int `yaml:"multis,omitempty"`
	Total  int `yaml:"total,omitempty"`
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
