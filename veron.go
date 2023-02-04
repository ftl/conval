/*
This file contains the implementation of the specific things for contests announced by the Dutch amateur radio association VERON.
*/
package conval

import (
	"regexp"
	"strings"

	"github.com/ftl/hamradio/callsign"
)

func init() {
	PropertyValidators[PAProvinceProperty] = RegexpValidator(validPAProvince, "PA province")

	PropertyGetters[PAProvinceProperty] = getTheirExchangeProperty(PAProvinceProperty)
	PropertyGetters[VeronEntityProperty] = PropertyGetterFunc(getVeronEntity)
}

const (
	PAProvinceProperty  Property = "pa_province"
	VeronEntityProperty Property = "veron_entity"
)

var (
	validPAProvince = regexp.MustCompile(`DR|FL|FR|GD|GR|LB|NB|NH|OV|UT|ZH|ZL`)
)

func getVeronEntity(qso QSO) string {
	return VeronEntity(qso.TheirCall, qso.TheirCountry)
}

func VeronEntity(call callsign.Callsign, dxccEntity DXCCEntity) string {
	dxccEntity = DXCCEntity(strings.ToUpper(string(dxccEntity)))
	switch dxccEntity {
	case "CE", "JA", "LU", "PY", "K", "VK", "ZS", "ZL":
		return string(dxccEntity) + veronCallAreaNumber(call)
	case "VE":
		return CanadianTerritory(call)
	case "UA9":
		return "UA" + veronCallAreaNumber(call)
	default:
		return string(dxccEntity)
	}
}

var veronNumberCallAreaExpression = regexp.MustCompile("[0-9]+")

func veronCallAreaNumber(call callsign.Callsign) string {
	var number string
	switch {
	case call.Prefix != "":
		number = veronNumberCallAreaExpression.FindString(call.Prefix)
	case call.Suffix != "":
		number = veronNumberCallAreaExpression.FindString(call.Suffix)
	default:
		number = veronNumberCallAreaExpression.FindString(call.BaseCall[1:])
	}
	if number == "" {
		number = "0"
	}
	if len(number) > 1 {
		number = number[1:]
	}
	return number
}

var canadianPrefixExpression = regexp.MustCompile("[CVX][A-Z][0-9]")

func CanadianTerritory(call callsign.Callsign) string {
	// according to https://hamwaves.com/map.ca/en/index.html
	callPrefix := ""
	if callPrefix == "" && call.Prefix != "" {
		callPrefix = canadianPrefixExpression.FindString(call.Prefix)
	}
	if callPrefix == "" {
		callPrefix = canadianPrefixExpression.FindString(call.BaseCall)
	}

	switch callPrefix {
	case "VO1", "VO2", "VY1", "VY2", "VY9", "VY0", "CY0", "CY9":
		return callPrefix
	default:
		return "VE" + veronCallAreaNumber(call)
	}
}
