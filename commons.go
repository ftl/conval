/*
This file contains the implementation of the common properties used for contest scoring.
*/
package conval

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/ftl/hamradio/callsign"
	"github.com/ftl/hamradio/locator"
)

const (
	TheirCallProperty        Property = "their_call"
	RSTProperty              Property = "rst"
	SerialNumberProperty     Property = "serial"
	MemberNumberProperty     Property = "member_number"
	NoMemberProperty         Property = "nm"
	CallsignProperty         Property = "callsign" // can be used as exchange, e.g. in the silent key memorial contests
	CQZoneProperty           Property = "cq_zone"
	ITUZoneProperty          Property = "itu_zone"
	DXCCEntityProperty       Property = "dxcc_entity"
	WorkingConditionProperty Property = "working_condition"
	NameProperty             Property = "name"
	AgeProperty              Property = "age"
	PowerProperty            Property = "power"
	ClassProperty            Property = "class"
	StateProvinceProperty    Property = "state_province"
	DXCCPrefixProperty       Property = "dxcc_prefix" // can be used as exchange, e.g. in the CWops contests
	WAEEntityProperty        Property = "wae_entity"
	WPXPrefixProperty        Property = "wpx_prefix"
	ContinentProperty        Property = "continent" // can be used as exchange, e.g. in the LABRE-DX contest
	LocatorProperty          Property = "locator"   // can be used as exchange
	DistanceProperty         Property = "distance"  // derived, needs the locator property

	GenericTextProperty   Property = "generic_text"
	GenericNumberProperty Property = "generic_number"
	EmptyProperty         Property = "empty"
)

func init() {
	commonPropertyValidators[RSTProperty] = RegexpValidator(validRST, "report")
	commonPropertyValidators[SerialNumberProperty] = RegexpValidator(validSerialNumber, "serial number")
	commonPropertyValidators[MemberNumberProperty] = RegexpValidator(validMemberNumber, "member number")
	commonPropertyValidators[NoMemberProperty] = RegexpValidator(validNoMember, "no member")
	commonPropertyValidators[CallsignProperty] = CallsignValidator
	commonPropertyValidators[CQZoneProperty] = NumberRangeValidator(1, 40, "CQ zone")
	commonPropertyValidators[ITUZoneProperty] = NumberRangeValidator(1, 90, "ITU zone")
	commonPropertyValidators[NameProperty] = RegexpValidator(validName, "name")
	commonPropertyValidators[AgeProperty] = RegexpValidator(validGenericNumber, "age")
	commonPropertyValidators[PowerProperty] = RegexpValidator(validPower, "power")
	commonPropertyValidators[ClassProperty] = RegexpValidator(validName, "class")
	commonPropertyValidators[StateProvinceProperty] = RegexpValidator(validStateProvince, "state or province")
	commonPropertyValidators[DXCCPrefixProperty] = DXCCPrefixValidator
	commonPropertyValidators[ContinentProperty] = RegexpValidator(validContinent, "continent")
	commonPropertyValidators[LocatorProperty] = RegexpValidator(validContinent, "locator")
	commonPropertyValidators[GenericTextProperty] = RegexpValidator(validGenericText, "generic text")
	commonPropertyValidators[GenericNumberProperty] = RegexpValidator(validGenericNumber, "generic number")
	commonPropertyValidators[EmptyProperty] = EmptyValidator

	commonPropertyGetters[TheirCallProperty] = PropertyGetterFunc(getTheirCall)
	commonPropertyGetters[RSTProperty] = getTheirExchangeProperty(RSTProperty)
	commonPropertyGetters[SerialNumberProperty] = getTheirExchangeProperty(SerialNumberProperty)
	commonPropertyGetters[MemberNumberProperty] = getTheirExchangeProperty(MemberNumberProperty)
	commonPropertyGetters[NoMemberProperty] = getTheirExchangeProperty(NoMemberProperty)
	commonPropertyGetters[CQZoneProperty] = PropertyGetterFunc(getCQZone)
	commonPropertyGetters[ITUZoneProperty] = PropertyGetterFunc(getITUZone)
	commonPropertyGetters[DXCCEntityProperty] = PropertyGetterFunc(getDXCCEntity)
	commonPropertyGetters[WorkingConditionProperty] = PropertyGetterFunc(getCallsignWorkingCondition)
	commonPropertyGetters[NameProperty] = getTheirExchangeProperty(NameProperty)
	commonPropertyGetters[AgeProperty] = getTheirExchangeProperty(AgeProperty)
	commonPropertyGetters[PowerProperty] = getTheirExchangeProperty(PowerProperty)
	commonPropertyGetters[ClassProperty] = getTheirExchangeProperty(ClassProperty)
	commonPropertyGetters[StateProvinceProperty] = getTheirExchangeProperty(StateProvinceProperty)
	commonPropertyGetters[DXCCPrefixProperty] = getTheirExchangeProperty(DXCCPrefixProperty)
	commonPropertyGetters[WAEEntityProperty] = PropertyGetterFunc(getWAEEntity)
	commonPropertyGetters[WPXPrefixProperty] = PropertyGetterFunc(getWPXPrefix)
	commonPropertyGetters[ContinentProperty] = PropertyGetterFunc(getContinent)
	commonPropertyGetters[LocatorProperty] = getTheirExchangeProperty(LocatorProperty)
	commonPropertyGetters[DistanceProperty] = PropertyGetterFunc(getDistance)
	commonPropertyGetters[GenericTextProperty] = getTheirExchangeProperty(GenericTextProperty)
	commonPropertyGetters[GenericNumberProperty] = getTheirExchangeProperty(GenericNumberProperty)
	commonPropertyGetters[EmptyProperty] = PropertyGetterFunc(getEmpty)

	myPropertyGetters[RSTProperty] = getMyExchangeProperty(RSTProperty)
	myPropertyGetters[SerialNumberProperty] = getMyExchangeProperty(SerialNumberProperty)
	myPropertyGetters[MemberNumberProperty] = getMyExchangeProperty(MemberNumberProperty)
	myPropertyGetters[NoMemberProperty] = getMyExchangeProperty(NoMemberProperty)
	myPropertyGetters[CQZoneProperty] = PropertyGetterFunc(getMyCQZone)
	myPropertyGetters[ITUZoneProperty] = PropertyGetterFunc(getMyITUZone)
	myPropertyGetters[DXCCEntityProperty] = PropertyGetterFunc(getMyDXCCEntity)
	myPropertyGetters[WorkingConditionProperty] = PropertyGetterFunc(getMyWorkingCondition)
	myPropertyGetters[NameProperty] = getMyExchangeProperty(NameProperty)
	myPropertyGetters[AgeProperty] = getMyExchangeProperty(AgeProperty)
	myPropertyGetters[PowerProperty] = getMyExchangeProperty(PowerProperty)
	myPropertyGetters[ClassProperty] = getMyExchangeProperty(ClassProperty)
	myPropertyGetters[StateProvinceProperty] = getMyExchangeProperty(StateProvinceProperty)
	myPropertyGetters[DXCCPrefixProperty] = getMyExchangeProperty(DXCCPrefixProperty)
	myPropertyGetters[WAEEntityProperty] = PropertyGetterFunc(getMyWAEEntity)
	myPropertyGetters[WPXPrefixProperty] = PropertyGetterFunc(getMyWPXPrefix)
	myPropertyGetters[ContinentProperty] = PropertyGetterFunc(getMyContinent)
	myPropertyGetters[LocatorProperty] = PropertyGetterFunc(getMyLocator)
	myPropertyGetters[GenericTextProperty] = getTheirExchangeProperty(GenericTextProperty)
	myPropertyGetters[GenericNumberProperty] = getTheirExchangeProperty(GenericNumberProperty)
	myPropertyGetters[EmptyProperty] = PropertyGetterFunc(getEmpty)
}

// Common Exchange Validators

func RegexpValidator(exp *regexp.Regexp, name string) PropertyValidator {
	return PropertyValidatorFunc(func(exchange string, prefixes PrefixDatabase) error {
		exchange = strings.ToUpper(strings.TrimSpace(exchange))
		value := exp.FindString(exchange)
		if len(value) == 0 || len(value) != len(exchange) {
			return fmt.Errorf("%s is not a valid %s", exchange, name)
		}
		return nil
	})
}

func NumberRangeValidator(from, to int, name string) PropertyValidator {
	return PropertyValidatorFunc(func(exchange string, _ PrefixDatabase) error {
		exchange = strings.ToUpper(strings.TrimSpace(exchange))
		value, err := strconv.Atoi(exchange)
		if err != nil {
			return fmt.Errorf("%s is not a valid %s: %w", exchange, name, err)
		}
		if value < from || value > to {
			return fmt.Errorf("%s is not a valid %s", exchange, name)
		}
		return nil
	})
}

var EmptyValidator = PropertyValidatorFunc(func(exchange string, _ PrefixDatabase) error {
	exchange = strings.TrimSpace(exchange)
	if exchange != "" {
		return fmt.Errorf("%s is not empty", exchange)
	}
	return nil
})

var (
	validRST          = regexp.MustCompile(`[1-5][1-9][1-9]*`)
	validSerialNumber = regexp.MustCompile(`\d+`)
	validMemberNumber = regexp.MustCompile(`\d+`)
	validNoMember     = regexp.MustCompile(`(NM)?`)
	validName         = regexp.MustCompile(`[A-Z]+`)
	validPower        = regexp.MustCompile(`[A-Z0-9]+`)
	// according to https://contests.arrl.org/contestmultipliers.php
	validStateProvince = regexp.MustCompile(`AB|BC|LB|MB|NB|NF|NS|NT|NU|ON|PE|QC|SK|YT|AL|AK|AZ|AR|CA|CO|CT|DC|DE|FL|GA|HI|ID|IL|IN|IA|KS|KY|LA|ME|MD|MA|MI|MN|MS|MO|MT|NE|NV|NH|NJ|NM|NY|NC|ND|OH|OK|OR|PA|RI|SC|SD|TN|TX|UT|VT|VA|WA|WV|WI|WY`)
	validContinent     = regexp.MustCompile(`AF|AN|AS|EU|OC|NA|SA`)
	validGenericText   = regexp.MustCompile(`[A-Z][A-Z0-9]*`)
	validGenericNumber = regexp.MustCompile(`[0-9]*`)

	CallsignValidator = PropertyValidatorFunc(func(exchange string, _ PrefixDatabase) error {
		_, err := callsign.Parse(exchange)
		if err != nil {
			return err
		}
		return nil
	})
	DXCCPrefixValidator = PropertyValidatorFunc(func(exchange string, prefixes PrefixDatabase) error {
		_, entity, _, _, found := prefixes.Find(exchange)
		if !found {
			return fmt.Errorf("%s is not a valid DXCC prefix", exchange)
		}
		if exchange != string(entity) {
			return fmt.Errorf("%s is not a primary DXCC prefix", exchange)
		}
		return nil
	})
	LocatorValidator = PropertyValidatorFunc(func(exchange string, _ PrefixDatabase) error {
		_, err := locator.Parse(exchange)
		if err != nil {
			return err
		}
		return nil
	})
)

// Common Property Getters

func getTheirCall(qso QSO, _ Setup, _ PrefixDatabase) string {
	return qso.TheirCall.String()
}

func getMyCQZone(qso QSO, setup Setup, prefixes PrefixDatabase) string {
	exchange, ok := qso.MyExchange[CQZoneProperty]
	if ok {
		return exchange
	}
	_, _, cqZone, _, ok := prefixes.Find(setup.MyCall.String())
	if ok {
		return cqZone.String()
	}
	return ""
}

func getCQZone(qso QSO, _ Setup, prefixes PrefixDatabase) string {
	exchange, ok := qso.TheirExchange[CQZoneProperty]
	if ok {
		return exchange
	}
	_, _, cqZone, _, ok := prefixes.Find(qso.TheirCall.String())
	if ok {
		return cqZone.String()
	}
	return ""
}

func getMyITUZone(qso QSO, setup Setup, prefixes PrefixDatabase) string {
	exchange, ok := qso.MyExchange[ITUZoneProperty]
	if ok {
		return exchange
	}
	_, _, _, ituZone, ok := prefixes.Find(setup.MyCall.String())
	if ok {
		return ituZone.String()
	}
	return ""
}

func getITUZone(qso QSO, _ Setup, prefixes PrefixDatabase) string {
	exchange, ok := qso.TheirExchange[ITUZoneProperty]
	if ok {
		return exchange
	}
	_, _, _, ituZone, ok := prefixes.Find(qso.TheirCall.String())
	if ok {
		return ituZone.String()
	}
	return ""
}

func getMyDXCCEntity(qso QSO, setup Setup, prefixes PrefixDatabase) string {
	if setup.MyCountry != "" {
		return string(setup.MyCountry)
	}
	_, entity, _, _, ok := prefixes.Find(setup.MyCall.String())
	if ok {
		return entity.String()
	}
	return ""
}

func getDXCCEntity(qso QSO, _ Setup, prefixes PrefixDatabase) string {
	if qso.TheirCountry != "" {
		return string(qso.TheirCountry)
	}
	_, entity, _, _, ok := prefixes.Find(qso.TheirCall.String())
	if ok {
		return entity.String()
	}
	return ""
}

func getMyWAEEntity(_ QSO, setup Setup, _ PrefixDatabase) string {
	return WAEEntity(setup.MyCall, setup.MyCountry)
}

func getWAEEntity(qso QSO, _ Setup, _ PrefixDatabase) string {
	return WAEEntity(qso.TheirCall, qso.TheirCountry)
}

func WAEEntity(call callsign.Callsign, dxccEntity DXCCEntity) string {
	dxccEntity = DXCCEntity(strings.ToUpper(string(dxccEntity)))
	switch dxccEntity {
	case "K", "VE", "VK", "ZL", "ZS", "JA", "BY", "PY":
		// special entities outside EU with numerical call areas
		return string(dxccEntity) + waeCallAreaNumber(call)
	case "UA9":
		// asian russia is even more special
		return "UA" + waeCallAreaNumber(call)
	default:
		return string(dxccEntity)
	}
}

var waeNumberCallAreaExpression = regexp.MustCompile("[0-9]+")

func waeCallAreaNumber(call callsign.Callsign) string {
	var number string
	if number == "" && call.Prefix != "" {
		number = waeNumberCallAreaExpression.FindString(call.Prefix)
	}
	if number == "" && call.Suffix != "" {
		number = waeNumberCallAreaExpression.FindString(call.Suffix)
	}
	if number == "" {
		number = waeNumberCallAreaExpression.FindString(call.BaseCall[1:])
	}
	if len(number) > 1 {
		number = number[1:]
	}
	return number
}

func IsWAECountry(call callsign.Callsign, dxccEntity DXCCEntity) bool {
	dxccEntity = DXCCEntity(strings.ToUpper(string(dxccEntity)))
	switch dxccEntity {
	case "IG9":
		return false

	// this is the WAE country list according to https://www.darc.de/der-club/referate/conteste/wae-dx-contest/en/wae-rules/
	case "1A", "3A", "4O", "4U1I", "4U1V", "9A", "9H", "C3", "CT", "CU", "DL",
		"E7", "EA", "EA6", "EI", "ER", "ES", "EU", "F", "G", "GD", "GI", "GJ",
		"GM", "GU", "GW", "HA", "HB", "HB0", "HV", "I", "IS", "IT9", "JW", "JW/b",
		"JX", "LA", "LX", "LY", "LZ", "OE", "OH", "OH0", "OJ0", "OK", "OM", "ON",
		"OY", "OZ", "PA", "UA", "UA2", "S5", "SM", "SP", "SV", "SV/a", "SV5", "SV9",
		"T7", "TA1", "TF", "TK", "UR", "YL", "YO", "YU", "Z6", "Z3", "ZA", "ZB":
		return true

	default:
		return false
	}
}

func getMyWorkingCondition(_ QSO, setup Setup, _ PrefixDatabase) string {
	return setup.MyCall.WorkingCondition
}

func getCallsignWorkingCondition(qso QSO, _ Setup, _ PrefixDatabase) string {
	return qso.TheirCall.WorkingCondition
}

func getMyExchangeProperty(property Property) PropertyGetter {
	return PropertyGetterFunc(func(qso QSO, _ Setup, _ PrefixDatabase) string {
		return qso.MyExchange[property]
	})
}

func getTheirExchangeProperty(property Property) PropertyGetter {
	return PropertyGetterFunc(func(qso QSO, _ Setup, _ PrefixDatabase) string {
		return qso.TheirExchange[property]
	})
}

func getEmpty(_ QSO, _ Setup, _ PrefixDatabase) string {
	return ""
}

func getMyWPXPrefix(_ QSO, setup Setup, _ PrefixDatabase) string {
	return WPXPrefix(setup.MyCall)
}

func getWPXPrefix(qso QSO, _ Setup, _ PrefixDatabase) string {
	return WPXPrefix(qso.TheirCall)
}

var parseWPXPrefixExpression = regexp.MustCompile(`^[A-Z0-9]?[A-Z][0-9]*`)

func WPXPrefix(call callsign.Callsign) string {
	var p string
	if p == "" && call.Prefix != "" {
		p = parseWPXPrefixExpression.FindString(call.Prefix)
	}
	if p == "" && call.Suffix != "" {
		p = parseWPXPrefixExpression.FindString(call.Suffix)
	}
	if p == "" {
		p = parseWPXPrefixExpression.FindString(call.BaseCall)
	}
	if p == "" {
		return ""
	}
	runes := []rune(p)
	if !unicode.IsDigit(runes[len(runes)-1]) {
		p += "0"
	}
	return p
}

func getMyContinent(_ QSO, setup Setup, prefixes PrefixDatabase) string {
	return string(ContinentForCallsign(setup.MyCall, prefixes))
}

func getContinent(qso QSO, _ Setup, prefixes PrefixDatabase) string {
	return string(ContinentForCallsign(qso.TheirCall, prefixes))
}

func ContinentForCallsign(call callsign.Callsign, prefixes PrefixDatabase) Continent {
	result, _, _, _, found := prefixes.Find(call.String())
	if !found {
		return ""
	}
	return result
}

func getMyLocator(_ QSO, setup Setup, _ PrefixDatabase) string {
	return setup.GridLocator.String()
}

func getDistance(qso QSO, setup Setup, _ PrefixDatabase) string {
	distance := Distance(setup.GridLocator, qso.TheirExchange[LocatorProperty])
	return strconv.Itoa(distance)
}

func Distance(myLocator locator.Locator, theirLocatorExchange string) int {
	if myLocator.IsZero() {
		return 0
	}
	theirLocator, err := locator.Parse(theirLocatorExchange)
	if err != nil {
		return 0
	}
	return int(locator.Distance(myLocator, theirLocator))
}
