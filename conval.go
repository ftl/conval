/*
The package conval helps to evaluate the log files from amateur radio contests.
*/
package conval

import (
	"time"

	"github.com/ftl/hamradio/callsign"
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
	ModeRTTY    Mode = "rtty"
	ModeDigital Mode = "digital"
)

type PropertyValidator interface {
	ValidateProperty(string) error
}

type PropertyValidatorFunc func(string) error

func (f PropertyValidatorFunc) ValidateProperty(exchange string) error {
	return f(exchange)
}

var PropertyValidators = map[Property]PropertyValidator{
	TheirRSTProperty:     PropertyValidatorFunc(ValidateRST),
	SerialNumberProperty: PropertyValidatorFunc(ValidateSerialNumber),
	MemberNumberProperty: PropertyValidatorFunc(ValidateMemberNumber),
	CQZoneProperty:       PropertyValidatorFunc(ValidateCQZone),
	ITUZoneProperty:      PropertyValidatorFunc(ValidateITUZone),
	NoMemberProperty:     PropertyValidatorFunc(ValidateNoMember),
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
)

type DXCCEntity string

const (
	SameCountry  DXCCEntity = "same"
	OtherCountry DXCCEntity = "other"
)

type BandRule string

const (
	OncePerBand BandRule = "once_per_band"
	Once        BandRule = "once"
)

type Overlay string

const (
	ClassicOverlay           Overlay = "classic"
	ThreeBandAndWiresOverlay Overlay = "tb_wires"
	RookieOverlay            Overlay = "rookie"
	YouthOverlay             Overlay = "youth"
)

type Property string

const (
	TheirRSTProperty     Property = "rst"
	SerialNumberProperty Property = "serial"
	MemberNumberProperty Property = "member_number"
	NoMemberProperty     Property = "nm"
	CQZoneProperty       Property = "cq_zone"
	ITUZoneProperty      Property = "itu_zone"
	DXCCEntityProperty   Property = "dxcc_entity"
	WPXPrefixProperty    Property = "wpx_prefix"
)

type PropertyGetter interface {
	GetProperty(QSO) string
}

type PropertyGetterFunc func(QSO) string

func (f PropertyGetterFunc) GetProperty(qso QSO) string {
	return f(qso)
}

var PropertyGetters = map[Property]PropertyGetter{
	TheirRSTProperty:     GetTheirExchangeProperty(TheirRSTProperty),
	SerialNumberProperty: GetTheirExchangeProperty(SerialNumberProperty),
	MemberNumberProperty: GetTheirExchangeProperty(MemberNumberProperty),
	NoMemberProperty:     GetTheirExchangeProperty(NoMemberProperty),
	CQZoneProperty:       PropertyGetterFunc(GetCQZone),
	ITUZoneProperty:      PropertyGetterFunc(GetITUZone),
	DXCCEntityProperty:   PropertyGetterFunc(GetDXCCEntity),
	WPXPrefixProperty:    PropertyGetterFunc(GetWPXPrefix),
}

type QSO struct {
	TheirCall      callsign.Callsign
	TheirContinent Continent
	TheirCountry   DXCCEntity

	Timestamp time.Time
	Band      ContestBand

	MyExchange    QSOExchange
	TheirExchange QSOExchange
}

type QSOExchange map[Property]string

type Setup struct {
	MyCall     callsign.Callsign
	Operator   OperatorMode
	Overlay    Overlay
	Bands      []ContestBand
	Modes      []Mode
	MyExchange QSOExchange

	MyContinent Continent
	MyCountry   DXCCEntity
}

type Constraint struct {
	Operator OperatorMode `yaml:"operator"`
	Overlay  Overlay      `yaml:"overlay"`
}
