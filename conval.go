/*
The package conval helps to evaluate the log files from amateur radio contests.
*/
package conval

import (
	"log"
	"strings"
	"time"

	"github.com/ftl/hamradio/callsign"
	"github.com/ftl/hamradio/dxcc"
	"github.com/ftl/hamradio/locator"
)

type OperatorMode string

const (
	SingleOperator OperatorMode = "single"
	MultiOperator  OperatorMode = "multi"
)

type TXMode string

const (
	OneTX         TXMode = "one"
	TwoTX         TXMode = "two"
	MultiTX       TXMode = "multi"
	DistributedTX TXMode = "distributed"
)

type PowerMode string

const (
	HighPower PowerMode = "high"
	LowPower  PowerMode = "low"
	QRPPower  PowerMode = "qrp"
)

type ContestBand string

const (
	BandAll ContestBand = "all"

	Band160m ContestBand = "160m"
	Band80m  ContestBand = "80m"
	Band40m  ContestBand = "40m"
	Band20m  ContestBand = "20m"
	Band15m  ContestBand = "15m"
	Band10m  ContestBand = "10m"
)

var AllHFBands = []ContestBand{Band160m, Band80m, Band40m, Band20m, Band15m, Band10m}

type BandMode string

const (
	AllBands   BandMode = "all"
	SingleBand BandMode = "single"
)

type Mode string

const (
	ModeALL Mode = "all"

	ModeCW      Mode = "cw"
	ModeSSB     Mode = "ssb"
	ModeFM      Mode = "fm"
	ModeRTTY    Mode = "rtty"
	ModeDigital Mode = "digital"
)

type PropertyValidator interface {
	ValidateProperty(string, PrefixDatabase) error
}

type PropertyValidatorFunc func(string, PrefixDatabase) error

func (f PropertyValidatorFunc) ValidateProperty(exchange string, prefixes PrefixDatabase) error {
	return f(exchange, prefixes)
}

var PropertyValidators = map[Property]PropertyValidator{}

type Continent string

const (
	Africa       Continent = "af"
	Antarctica   Continent = "an"
	Asia         Continent = "as"
	Europa       Continent = "eu"
	NorthAmerica Continent = "na"
	Oceania      Continent = "oc"
	SouthAmerica Continent = "sa"

	SameContinent  Continent = "same"
	OtherContinent Continent = "other"
	NotContinent   Continent = "not"
)

type DXCCEntity string

const (
	SameCountry  DXCCEntity = "same"
	OtherCountry DXCCEntity = "other"
	NotCountry   DXCCEntity = "not"
)

type PrefixDatabase interface {
	Find(s string) (Continent, DXCCEntity, bool)
}

func NewPrefixDatabase() (*prefixDatabase, error) {
	prefixes, _, err := dxcc.DefaultPrefixes(true)
	if err != nil {
		return nil, err
	}
	return &prefixDatabase{prefixes}, nil
}

type prefixDatabase struct {
	prefixes *dxcc.Prefixes
}

func (d prefixDatabase) Find(s string) (Continent, DXCCEntity, bool) {
	entities, found := d.prefixes.Find(s)
	if !found || len(entities) == 0 {
		return "", "", false
	}

	return Continent(strings.ToLower(entities[0].Continent)), DXCCEntity(strings.ToLower(entities[0].PrimaryPrefix)), true
}

type BandRule string

const (
	Once               BandRule = "once"
	OncePerBand        BandRule = "once_per_band"
	OncePerBandAndMode BandRule = "once_per_band_and_mode"
)

type Overlay string

const (
	NoOverlay                Overlay = ""
	ClassicOverlay           Overlay = "classic"
	ThreeBandAndWiresOverlay Overlay = "tb_wires"
	RookieOverlay            Overlay = "rookie"
	YouthOverlay             Overlay = "youth"
)

type Property string

type PropertyGetter interface {
	GetProperty(QSO) string
}

type PropertyGetterFunc func(QSO) string

func (f PropertyGetterFunc) GetProperty(qso QSO) string {
	return f(qso)
}

var PropertyGetters = map[Property]PropertyGetter{}

type QSO struct {
	TheirCall      callsign.Callsign
	TheirContinent Continent
	TheirCountry   DXCCEntity

	Timestamp time.Time
	Band      ContestBand
	Mode      Mode

	MyExchange    QSOExchange
	TheirExchange QSOExchange
}

type QSOExchange map[Property]string

func ParseExchange(fields []ExchangeField, values []string, prefixes PrefixDatabase) QSOExchange {
	result := make(QSOExchange)
	result.Add(fields, values, prefixes)
	return result
}

func (e QSOExchange) Add(fields []ExchangeField, values []string, prefixes PrefixDatabase) {
	for i, field := range fields {
		if i >= len(values) {
			break
		}
		value := strings.ToUpper(strings.TrimSpace(values[i]))
		for _, property := range field {
			validator, ok := PropertyValidators[property]
			if !ok {
				log.Printf("no validator for property %s", property)
				continue
			}
			err := validator.ValidateProperty(value, prefixes)
			if err == nil {
				e[property] = value
				break
			}
		}
	}
}

type Setup struct {
	MyCall      callsign.Callsign
	MyContinent Continent
	MyCountry   DXCCEntity

	GridLocator locator.Locator
	Operators   []callsign.Callsign

	OperatorMode OperatorMode
	Overlay      Overlay
	Power        PowerMode
	Bands        []ContestBand
	Modes        []Mode

	MyExchange QSOExchange
}

type Constraint struct {
	OperatorMode OperatorMode `yaml:"operator_mode"`
	Overlay      Overlay      `yaml:"overlay"`
}
