/*
This file contains the implementation of the specific things for contests announced by the DARC.
*/
package conval

import (
	"regexp"
	"strings"

	"github.com/ftl/hamradio/callsign"
)

func init() {
	PropertyValidators[WAGDOKProperty] = RegexpValidator(validDOK, "DOK")

	PropertyGetters[WAEEntityProperty] = PropertyGetterFunc(getWAEEntity)
	PropertyGetters[WAGDistrictProperty] = PropertyGetterFunc(getWAGDistrict)
}

const (
	WAGDOKProperty      Property = "wag_dok"
	WAGDistrictProperty Property = "wag_district"
	WAEEntityProperty   Property = "wae_entity"
)

var validDOK = regexp.MustCompile(`\d*[A-Z][A-Z0-9ÄÖÜ-]*`)

func getWAGDistrict(qso QSO) string {
	dok, ok := qso.TheirExchange[WAGDOKProperty]
	if !ok {
		return ""
	}
	dok = strings.ToUpper(strings.TrimSpace(dok))
	if len(dok) == 0 {
		return ""
	}
	isLetter := func(b byte) bool {
		return b >= 'A' && b <= 'Z'
	}
	for i := range dok {
		if isLetter(dok[i]) {
			return string(dok[i])
		}
	}
	return ""
}

func getWAEEntity(qso QSO) string {
	return WAEEntity(qso.TheirCall, qso.TheirCountry)
}

func WAEEntity(call callsign.Callsign, dxccEntity DXCCEntity) string {
	dxccEntity = DXCCEntity(strings.ToUpper(string(dxccEntity)))
	switch dxccEntity {
	case "K", "VE", "VK", "ZL", "ZS", "JA", "BY", "PY":
		// special entities outside EU with numerical call areas
		return string(dxccEntity) + waeCallAreaNumber(call, dxccEntity)
	case "UA9":
		// asian russia is even more special
		return "UA" + waeCallAreaNumber(call, dxccEntity)
	default:
		return string(dxccEntity)
	}
}

var waeNumberCallAreaExpression = regexp.MustCompile("[0-9]+")

func waeCallAreaNumber(call callsign.Callsign, dxccEntity DXCCEntity) string {
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
