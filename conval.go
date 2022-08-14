/*
The package conval helps to evaluate the log files from amateur radio contests.
*/
package conval

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
	ModeCW      Mode = "cw"
	ModeSSB     Mode = "ssb"
	ModeRTTY    Mode = "rtty"
	ModeDigital Mode = "digital"
)

type Exchange string

const (
	RSTExchange          Exchange = "rst"
	SerialExchange       Exchange = "serial"
	MemberNumberExchange Exchange = "member_number"
	CQZoneExchange       Exchange = "cq_zone"
	ITUZoneExchange      Exchange = "itu_zone"
	NoMemberExchange     Exchange = "nm"

	RefDepartmentExchange Exchange = "ref_department" // TODO move into a specific REF package
	WAGDOKExchange        Exchange = "wag_dok"        // TODO move into a specific DARC package
)

type ExchangeValidator interface {
	ValidateExchange(string) error
}

type ExchangeValidatorFunc func(string) error

func (f ExchangeValidatorFunc) ValidateExchange(exchange string) error {
	return f(exchange)
}

var ExchangeValidators = map[Exchange]ExchangeValidator{
	RSTExchange: ExchangeValidatorFunc(func(exchange string) error {
		return nil
	}), // TODO extract function into a separate file for common validators
}

type Continent string

const (
	Africa       Continent = "af"
	Antarctica   Continent = "ar"
	Asia         Continent = "as"
	Europa       Continent = "eu"
	NorthAmerica Continent = "na"
	Oceania      Continent = "oc"
	SouthAmerica Continent = "sa"

	SameContinent  Continent = "same"
	OtherContinent Continent = "other"
)

type BandRule string

const (
	OncePerBand BandRule = "once_per_band"
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
	CQZoneProperty     Property = "cq_zone"
	ITUZoneProperty    Property = "itu_zone"
	DXCCEntityProperty Property = "dxcc_entity"
	WPXPrefixProperty  Property = "wpx_prefix"

	REFDepartmentProperty Property = "ref_department" // TODO move into a specific REF package
	WAEEntityProperty     Property = "wae_property"   // TODO move into a specific DARC package
	WAGDistrictProperty   Property = "wag_district"   // TODO move into a specific DARC package
)

type PropertyGetter interface {
	GetProperty(QSO) string
}

type PropertyGetterFunc func(QSO) string

func (f PropertyGetterFunc) GetProperty(qso QSO) string {
	return f(qso)
}

var PropertyGetters = map[Property]PropertyGetter{
	CQZoneProperty: PropertyGetterFunc(func(QSO) string {
		return ""
	}), // TODO extract function into a separate file for common property getters
}

type QSO struct {
	// TODO define
}

type Setup struct {
	// TODO define
}

type Constraint struct {
	// TODO define
}
