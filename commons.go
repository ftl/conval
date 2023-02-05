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

	GenericTextProperty   Property = "generic_text"
	GenericNumberProperty Property = "generic_number"
	EmptyProperty         Property = "empty"
)

func init() {
	PropertyValidators[RSTProperty] = RegexpValidator(validRST, "report")
	PropertyValidators[SerialNumberProperty] = RegexpValidator(validSerialNumber, "serial number")
	PropertyValidators[MemberNumberProperty] = RegexpValidator(validMemberNumber, "member number")
	PropertyValidators[NoMemberProperty] = RegexpValidator(validNoMember, "no member")
	PropertyValidators[CallsignProperty] = CallsignValidator
	PropertyValidators[CQZoneProperty] = NumberRangeValidator(1, 40, "CQ zone")
	PropertyValidators[ITUZoneProperty] = NumberRangeValidator(1, 90, "ITU zone")
	PropertyValidators[NameProperty] = RegexpValidator(validName, "name")
	PropertyValidators[AgeProperty] = RegexpValidator(validGenericNumber, "age")
	PropertyValidators[PowerProperty] = RegexpValidator(validPower, "power")
	PropertyValidators[ClassProperty] = RegexpValidator(validName, "class")
	PropertyValidators[StateProvinceProperty] = RegexpValidator(validStateProvince, "state or province")
	PropertyValidators[DXCCPrefixProperty] = DXCCPrefixValidator
	PropertyValidators[GenericTextProperty] = RegexpValidator(validGenericText, "generic text")
	PropertyValidators[GenericNumberProperty] = RegexpValidator(validGenericNumber, "generic number")
	PropertyValidators[EmptyProperty] = EmptyValidator

	PropertyGetters[TheirCallProperty] = PropertyGetterFunc(getTheirCall)
	PropertyGetters[RSTProperty] = getTheirExchangeProperty(RSTProperty)
	PropertyGetters[SerialNumberProperty] = getTheirExchangeProperty(SerialNumberProperty)
	PropertyGetters[MemberNumberProperty] = getTheirExchangeProperty(MemberNumberProperty)
	PropertyGetters[NoMemberProperty] = getTheirExchangeProperty(NoMemberProperty)
	PropertyGetters[CQZoneProperty] = PropertyGetterFunc(getCQZone)
	PropertyGetters[ITUZoneProperty] = PropertyGetterFunc(getITUZone)
	PropertyGetters[DXCCEntityProperty] = PropertyGetterFunc(getDXCCEntity)
	PropertyGetters[WorkingConditionProperty] = PropertyGetterFunc(getCallsignWorkingCondition)
	PropertyGetters[NameProperty] = getTheirExchangeProperty(NameProperty)
	PropertyGetters[AgeProperty] = getTheirExchangeProperty(AgeProperty)
	PropertyGetters[PowerProperty] = getTheirExchangeProperty(PowerProperty)
	PropertyGetters[ClassProperty] = getTheirExchangeProperty(ClassProperty)
	PropertyGetters[StateProvinceProperty] = getTheirExchangeProperty(StateProvinceProperty)
	PropertyGetters[DXCCPrefixProperty] = getTheirExchangeProperty(DXCCPrefixProperty)
	PropertyGetters[GenericTextProperty] = getTheirExchangeProperty(GenericTextProperty)
	PropertyGetters[GenericNumberProperty] = getTheirExchangeProperty(GenericNumberProperty)
	PropertyGetters[EmptyProperty] = PropertyGetterFunc(getEmpty)
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
		_, entity, found := prefixes.Find(exchange)
		if !found {
			return fmt.Errorf("%s is not a valid DXCC prefix", exchange)
		}
		if exchange != string(entity) {
			return fmt.Errorf("%s is not a primary DXCC prefix", exchange)
		}
		return nil
	})
)

// Common Property Getters

func getTheirCall(qso QSO) string {
	return qso.TheirCall.String()
}

func getCQZone(qso QSO) string {
	exchange, ok := qso.TheirExchange[CQZoneProperty]
	if ok {
		return exchange
	}
	// TODO get CQ zone from a database
	return ""
}

func getITUZone(qso QSO) string {
	exchange, ok := qso.TheirExchange[ITUZoneProperty]
	if ok {
		return exchange
	}
	// TODO get ITU zone from a database
	return ""
}

func getDXCCEntity(qso QSO) string {
	if qso.TheirCountry != "" {
		return string(qso.TheirCountry)
	}
	// TODO get DXCC entity from a database
	return ""
}

func getCallsignWorkingCondition(qso QSO) string {
	return qso.TheirCall.WorkingCondition
}

func getTheirExchangeProperty(property Property) PropertyGetter {
	return PropertyGetterFunc(func(qso QSO) string {
		return qso.TheirExchange[property]
	})
}

func getEmpty(_ QSO) string {
	return ""
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
