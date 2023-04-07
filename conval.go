/*
The package conval helps to evaluate the log files from amateur radio contests.
*/
package conval

import (
	"log"
	"strconv"
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

type BandCount string

const (
	AllBands   BandCount = "all"
	SingleBand BandCount = "single"
)

func (c BandCount) ToInt() int {
	switch c {
	case AllBands:
		return 0
	case SingleBand:
		return 1
	default:
		i, err := strconv.Atoi(string(c))
		if err != nil {
			return 0
		}
		return i
	}
}

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

type PropertyValidators interface {
	PropertyValidator(Property) (PropertyValidator, bool)
}

type PropertyValidatorsFunc func(Property) (PropertyValidator, bool)

func (f PropertyValidatorsFunc) PropertyValidator(property Property) (PropertyValidator, bool) {
	return f(property)
}

var commonPropertyValidators = map[Property]PropertyValidator{}

func CommonPropertyValidator(property Property) (PropertyValidator, bool) {
	validator, ok := commonPropertyValidators[property]
	return validator, ok
}

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

func (e DXCCEntity) String() string {
	return string(e)
}

const (
	SameCountry  DXCCEntity = "same"
	OtherCountry DXCCEntity = "other"
	NotCountry   DXCCEntity = "not"
)

const (
	SamePrefix  string = "same"
	OtherPrefix string = "other"
	NotPrefix   string = "not"
)

type CQZone int

func (z CQZone) String() string {
	return strconv.Itoa(int(z))
}

type ITUZone int

func (z ITUZone) String() string {
	return strconv.Itoa(int(z))
}

type PrefixDatabase interface {
	Find(s string) (Continent, DXCCEntity, CQZone, ITUZone, bool)
}

type PrefixDatabaseFunc func(s string) (Continent, DXCCEntity, bool)

func (f PrefixDatabaseFunc) Find(s string) (Continent, DXCCEntity, bool) {
	return f(s)
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

func (d prefixDatabase) Find(s string) (Continent, DXCCEntity, CQZone, ITUZone, bool) {
	entities, found := d.prefixes.Find(s)
	if !found || len(entities) == 0 {
		return "", "", 0, 0, false
	}

	return Continent(strings.ToLower(entities[0].Continent)),
		DXCCEntity(strings.ToLower(entities[0].PrimaryPrefix)),
		CQZone(entities[0].CQZone),
		ITUZone(entities[0].ITUZone),
		true
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
	GetProperty(QSO, PrefixDatabase) string
}

type PropertyGetterFunc func(QSO, PrefixDatabase) string

func (f PropertyGetterFunc) GetProperty(qso QSO, prefixes PrefixDatabase) string {
	return f(qso, prefixes)
}

var commonPropertyGetters = map[Property]PropertyGetter{}

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

func (q QSO) TheirPrefix() string {
	return WPXPrefix(q.TheirCall)
}

type QSOExchange map[Property]string

func ParseExchange(fields []ExchangeField, values []string, prefixes PrefixDatabase, propertyValidators PropertyValidators) QSOExchange {
	result := make(QSOExchange)
	result.Add(fields, values, prefixes, propertyValidators)
	return result
}

func (e QSOExchange) Add(fields []ExchangeField, values []string, prefixes PrefixDatabase, propertyValidators PropertyValidators) {
	for i, field := range fields {
		if i >= len(values) {
			break
		}
		value := strings.ToUpper(strings.TrimSpace(values[i]))
		for _, property := range field {
			validator, ok := propertyValidators.PropertyValidator(property)
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

func (s Setup) MyPrefix() string {
	return WPXPrefix(s.MyCall)
}

type Constraint struct {
	OperatorMode OperatorMode `yaml:"operator_mode"`
	Overlay      Overlay      `yaml:"overlay"`
}
