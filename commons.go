/*
This file contains the implementation of the common properties used for contest scoring.
*/
package conval

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ftl/hamradio/callsign"
)

const (
	TheirCallProperty        Property = "their_call"
	TheirRSTProperty         Property = "rst"
	SerialNumberProperty     Property = "serial"
	MemberNumberProperty     Property = "member_number"
	NoMemberProperty         Property = "nm"
	CallsignProperty         Property = "callsign" // can be used as exchange, e.g. in the silent key memorial contests
	CQZoneProperty           Property = "cq_zone"
	ITUZoneProperty          Property = "itu_zone"
	DXCCEntityProperty       Property = "dxcc_entity"
	WorkingConditionProperty Property = "working_condition"
	NameProperty             Property = "name"
	StateProvinceProperty    Property = "state_province"
	DXCCPrefixProperty       Property = "dxcc_prefix" // can be used as exchange, e.g. in the CWops contests
	AlphanumProperty         Property = "alphanum"
)

func init() {
	PropertyValidators[TheirRSTProperty] = RegexpValidator(validRST, "report")
	PropertyValidators[SerialNumberProperty] = RegexpValidator(validSerialNumber, "serial number")
	PropertyValidators[MemberNumberProperty] = RegexpValidator(validMemberNumber, "member number")
	PropertyValidators[NoMemberProperty] = RegexpValidator(validNoMember, "no member")
	PropertyValidators[CallsignProperty] = CallsignValidator
	PropertyValidators[CQZoneProperty] = NumberRangeValidator(1, 40, "CQ zone")
	PropertyValidators[ITUZoneProperty] = NumberRangeValidator(1, 90, "ITU zone")
	PropertyValidators[NameProperty] = RegexpValidator(validName, "name")
	PropertyValidators[StateProvinceProperty] = RegexpValidator(validStateProvince, "state or province")
	PropertyValidators[DXCCPrefixProperty] = DXCCPrefixValidator
	PropertyValidators[AlphanumProperty] = RegexpValidator(validAlphanum, "alpha numeric")

	PropertyGetters[TheirCallProperty] = PropertyGetterFunc(getTheirCall)
	PropertyGetters[TheirRSTProperty] = getTheirExchangeProperty(TheirRSTProperty)
	PropertyGetters[SerialNumberProperty] = getTheirExchangeProperty(SerialNumberProperty)
	PropertyGetters[MemberNumberProperty] = getTheirExchangeProperty(MemberNumberProperty)
	PropertyGetters[NoMemberProperty] = getTheirExchangeProperty(NoMemberProperty)
	PropertyGetters[CQZoneProperty] = PropertyGetterFunc(getCQZone)
	PropertyGetters[ITUZoneProperty] = PropertyGetterFunc(getITUZone)
	PropertyGetters[DXCCEntityProperty] = PropertyGetterFunc(getDXCCEntity)
	PropertyGetters[WorkingConditionProperty] = PropertyGetterFunc(getCallsignWorkingCondition)
	PropertyGetters[NameProperty] = getTheirExchangeProperty(NameProperty)
	PropertyGetters[StateProvinceProperty] = getTheirExchangeProperty(StateProvinceProperty)
	PropertyGetters[DXCCPrefixProperty] = getTheirExchangeProperty(DXCCPrefixProperty)
	PropertyGetters[AlphanumProperty] = getTheirExchangeProperty(AlphanumProperty)
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

var (
	validRST           = regexp.MustCompile(`[1-5][1-9][1-9]*`)
	validSerialNumber  = regexp.MustCompile(`\d+`)
	validMemberNumber  = regexp.MustCompile(`\d+`)
	validNoMember      = regexp.MustCompile(`(NM)?`)
	validName          = regexp.MustCompile(`[A-Z]+`)
	validStateProvince = regexp.MustCompile(`AB|AL|AK|AZ|AR|BC|CA|CO|CT|DE|FL|GA|HI|ID|IL|IN|IA|KS|KY|LA|ME|MD|MA|MB|MI|MN|MS|MO|MT|NB|NE|NV|NH|NJ|NL|NM|NS|NY|NC|ND|OH|OK|ON|OR|PA|PE|QC|RI|SC|SD|SK|TN|TX|UT|VT|VA|WA|WV|WI|WY`)
	validAlphanum      = regexp.MustCompile(`[A-Z][A-Z0-9]*`)

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
