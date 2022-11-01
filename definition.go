package conval

import (
	"io"
	"time"

	"github.com/ftl/hamradio/callsign"
	"gopkg.in/yaml.v3"
)

type Definition struct {
	Name            string                `yaml:"name"`
	Identifier      string                `yaml:"identifier"` // see https://www.contestcalendar.com/cabnames.php
	OfficialRules   string                `yaml:"official_rules"`
	UploadURL       string                `yaml:"upload_url"`
	UploadFormat    string                `yaml:"upload_format"`
	Durations       []ConstrainedDuration `yaml:"durations"`
	Breaks          []ConstrainedDuration `yaml:"breaks"`
	Categories      []Category            `yaml:"categories"`
	Overlays        []Overlay             `yaml:"overlays"`
	Modes           []Mode                `yaml:"modes"`
	Bands           []ContestBand         `yaml:"bands"`
	BandChangeRules []BandChangeRule      `yaml:"band_change_rules"`
	Exchange        []ExchangeDefinition  `yaml:"exchange"`
	Scoring         Scoring               `yaml:"scoring"`
	Examples        []Example             `yaml:"examples"`
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
	Operator OperatorMode `yaml:"operator"`
	TX       TXMode       `yaml:"tx"`
	Power    PowerMode    `yaml:"power"`
	Bands    BandMode     `yaml:"bands"`
	Modes    []Mode       `yaml:"modes"`
	Assisted bool         `yaml:"assisted"`
}

type ExchangeDefinition struct {
	MyContinent           []Continent     `yaml:"my_continent"`
	MyCountry             []DXCCEntity    `yaml:"my_country"`
	TheirContinent        []Continent     `yaml:"their_continent"`
	TheirCountry          []DXCCEntity    `yaml:"their_country"`
	TheirWorkingCondition string          `yaml:"their_working_condition"`
	AdditionalWeight      int             `yaml:"additional_weight"`
	Fields                []ExchangeField `yaml:"fields"`
}

type ExchangeField []Property

type Scoring struct {
	QSORules    []ScoringRule `yaml:"qsos"`
	QSOBandRule BandRule      `yaml:"qso_band_rule"`
	MultiRules  []ScoringRule `yaml:"multis"`
}

type ScoringRule struct {
	MyContinent           []Continent   `yaml:"my_continent"`
	MyCountry             []DXCCEntity  `yaml:"my_country"`
	TheirContinent        []Continent   `yaml:"their_continent"`
	TheirCountry          []DXCCEntity  `yaml:"their_country"`
	TheirWorkingCondition string        `yaml:"their_working_condition"`
	Bands                 []ContestBand `yaml:"bands"`
	Property              Property      `yaml:"property"`
	BandRule              BandRule      `yaml:"band_rule"`
	AdditionalWeight      int           `yaml:"additional_weight"`
	Value                 int           `yaml:"value"`
}

type Example struct {
	Setup SetupExample `yaml:"setup"`
	QSOs  []QSOExample `yaml:"qsos"`
	Score BandScore    `yaml:"score"`
}

type SetupExample struct {
	MyCall      string     `yaml:"my_call"`
	MyContinent Continent  `yaml:"my_continent"`
	MyCountry   DXCCEntity `yaml:"my_country"`

	Operator OperatorMode  `yaml:"operator"`
	Overlay  Overlay       `yaml:"overlay"`
	Bands    []ContestBand `yaml:"bands"`
	Modes    []Mode        `yaml:"modes"`

	MyExchange QSOExchange `yaml:"my_exchange"`
}

func (s SetupExample) ToSetup() Setup {
	myCall, err := callsign.Parse(s.MyCall)
	if err != nil {
		myCall = callsign.Callsign{}
	}

	return Setup{
		MyCall:      myCall,
		MyContinent: s.MyContinent,
		MyCountry:   s.MyCountry,
		Operator:    s.Operator,
		Overlay:     s.Overlay,
		Bands:       s.Bands,
		Modes:       s.Modes,
		MyExchange:  s.MyExchange,
	}
}

type QSOExample struct {
	TheirCall      string     `yaml:"their_call"`
	TheirContinent Continent  `yaml:"their_continent"`
	TheirCountry   DXCCEntity `yaml:"their_country"`

	Timestamp time.Time   `yaml:"time"`
	Band      ContestBand `yaml:"band"`
	Mode      Mode        `yaml:"mode"`

	MyExchange    QSOExchange `yaml:"my_exchange"`
	TheirExchange []string    `yaml:"their_exchange"`

	Score QSOScore `yaml:",inline"`
}

func (q QSOExample) ToQSO(fields []ExchangeField) QSO {
	return QSO{
		TheirCall:      callsign.MustParse(q.TheirCall),
		TheirContinent: q.TheirContinent,
		TheirCountry:   q.TheirCountry,
		Timestamp:      q.Timestamp,
		Band:           q.Band,
		Mode:           q.Mode,
		MyExchange:     q.MyExchange,
		TheirExchange:  ParseExchange(fields, q.TheirExchange),
	}
}

func LoadYAML(r io.Reader) (*Definition, error) {
	decoder := yaml.NewDecoder(r)

	var result Definition
	err := decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
