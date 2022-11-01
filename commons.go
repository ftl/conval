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
	TheirRSTProperty         Property = "rst"
	SerialNumberProperty     Property = "serial"
	MemberNumberProperty     Property = "member_number"
	NoMemberProperty         Property = "nm"
	CallsignProperty         Property = "callsign" // can be used as exchanges, e.g. in the silent key memorial contests
	CQZoneProperty           Property = "cq_zone"
	ITUZoneProperty          Property = "itu_zone"
	DXCCEntityProperty       Property = "dxcc_entity"
	WorkingConditionProperty Property = "working_condition"
)

func init() {
	PropertyValidators[TheirRSTProperty] = RegexpValidator(validRST, "report")
	PropertyValidators[SerialNumberProperty] = RegexpValidator(validSerialNumber, "serial number")
	PropertyValidators[MemberNumberProperty] = RegexpValidator(validMemberNumber, "member number")
	PropertyValidators[NoMemberProperty] = RegexpValidator(validNoMember, "no member")
	PropertyValidators[CallsignProperty] = CallsignValidator
	PropertyValidators[CQZoneProperty] = NumberRangeValidator(1, 40, "CQ zone")
	PropertyValidators[ITUZoneProperty] = NumberRangeValidator(1, 90, "ITU zone")

	PropertyGetters[TheirRSTProperty] = getTheirExchangeProperty(TheirRSTProperty)
	PropertyGetters[SerialNumberProperty] = getTheirExchangeProperty(SerialNumberProperty)
	PropertyGetters[MemberNumberProperty] = getTheirExchangeProperty(MemberNumberProperty)
	PropertyGetters[NoMemberProperty] = getTheirExchangeProperty(NoMemberProperty)
	PropertyGetters[CQZoneProperty] = PropertyGetterFunc(getCQZone)
	PropertyGetters[ITUZoneProperty] = PropertyGetterFunc(getITUZone)
	PropertyGetters[DXCCEntityProperty] = PropertyGetterFunc(getDXCCEntity)
	PropertyGetters[WorkingConditionProperty] = PropertyGetterFunc(getCallsignWorkingCondition)
}

// Common Exchange Validators

func RegexpValidator(exp *regexp.Regexp, name string) PropertyValidator {
	return PropertyValidatorFunc(func(exchange string) error {
		exchange = strings.ToUpper(strings.TrimSpace(exchange))
		value := exp.FindString(exchange)
		if len(value) == 0 || len(value) != len(exchange) {
			return fmt.Errorf("%s is not a valid %s", exchange, name)
		}
		return nil
	})
}

func NumberRangeValidator(from, to int, name string) PropertyValidator {
	return PropertyValidatorFunc(func(exchange string) error {
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
	validRST          = regexp.MustCompile(`[1-5][1-9][1-9]*`)
	validSerialNumber = regexp.MustCompile(`\d+`)
	validMemberNumber = regexp.MustCompile(`\d+`)
	validNoMember     = regexp.MustCompile(`(NM)?`)

	CallsignValidator = PropertyValidatorFunc(func(exchange string) error {
		_, err := callsign.Parse(exchange)
		if err != nil {
			return err
		}
		return nil
	})
)

// Common Property Getters

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
