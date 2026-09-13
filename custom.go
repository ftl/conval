/*
This file contains the implementation of the contest specific properties:
- EUDXCC (https://eudxcc.altervista.org/)
- VERON, the Dutch amateur radio association (https://www.veron.nl)
- SAC, the Scandinavian Activity Contest (https://www.sactest.net)
- SARL, the South African Radio League (https://mysarl.org.za)
- BFRA, the Bulgarian Federation of Radio Amateurs (https://bfra.bg)
*/
package conval

import (
	"regexp"
	"strings"

	"github.com/ftl/hamradio/callsign"
)

func init() {
	commonPropertyValidators[EURegionProperty] = RegexpValidator(validEURegion, "EU region")
	commonPropertyValidators[PAProvinceProperty] = RegexpValidator(validPAProvince, "PA province")

	commonPropertyGetters[EURegionProperty] = getTheirExchangeProperty(EURegionProperty)
	commonPropertyGetters[PAProvinceProperty] = getTheirExchangeProperty(PAProvinceProperty)
	commonPropertyGetters[VeronEntityProperty] = PropertyGetterFunc(getVeronEntity)
	commonPropertyGetters[SACAreaProperty] = PropertyGetterFunc(getSACArea)
	commonPropertyGetters[SARLAreaProperty] = PropertyGetterFunc(getSARLArea)
	commonPropertyGetters[BalkanPrefixProperty] = PropertyGetterFunc(getBalkanPrefix)

	myPropertyGetters[PAProvinceProperty] = getMyExchangeProperty(PAProvinceProperty)
	myPropertyGetters[VeronEntityProperty] = PropertyGetterFunc(getMyVeronEntity)
}

const (
	EURegionProperty     Property = "eu_region"
	PAProvinceProperty   Property = "pa_province"
	VeronEntityProperty  Property = "veron_entity"
	SACAreaProperty      Property = "sac_area"
	SARLAreaProperty     Property = "sarl_area"
	BalkanPrefixProperty Property = "balkan_prefix"
)

var (
	validEURegion   = regexp.MustCompile(`AT0[1-9]|BE[01][0-9]|BG0[1-6]|CZ[01][0-9]|CY0[1-5]|DK0[1-6]|EE0[1-5]|FI[01][0-9]|FR[0-2][0-9]|DE[01][0-9]|GR[01][0-9]|HU0[1-7]|IE0[1-4]|IT[0-2][0-9]|LV0[1-6]|LT0[1-5]|LX01|MT0[1-5]|NL[01][0-9]|PL[01][0-9]|RO0[1-8]|SK0[1-8]|SI0[1-6]|ES[01][1-9]|SE[0-2][1-9]`)
	validPAProvince = regexp.MustCompile(`DR|FL|FR|GD|GR|LB|NB|NH|OV|UT|ZH|ZL`)
)

var callAreaDigitExpression = regexp.MustCompile("[0-9]")

func callAreaDigit(call callsign.Callsign) string {
	prefix := call.Prefix
	if prefix == "" {
		prefix = call.BaseCall
	}
	if len(prefix) > 1 {
		// the first character may be a digit that belongs to the prefix, not to the call area
		prefix = prefix[1:]
	}

	return callAreaDigitExpression.FindString(prefix)
}

func getMyVeronEntity(_ QSO, setup Setup, _ PrefixDatabase) string {
	return VeronEntity(setup.MyCall, setup.MyCountry)
}

func getVeronEntity(qso QSO, _ Setup, _ PrefixDatabase) string {
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

func getSACArea(qso QSO, setup Setup, prefixes PrefixDatabase) string {
	return sacArea(qso.TheirCall, DXCCEntity(getDXCCEntity(qso, setup, prefixes)))
}

func sacArea(call callsign.Callsign, dxccEntity DXCCEntity) string {
	// according to https://www.sactest.net/blog/scandinavian-activity-contest-2025-rules/ §8.2
	entity := strings.ToLower(string(dxccEntity))
	if entity == "" {
		return ""
	}

	number := callAreaDigit(call)
	if number == "" {
		number = "0"
	}
	return entity + number
}

func getSARLArea(qso QSO, setup Setup, prefixes PrefixDatabase) string {
	return sarlArea(qso.TheirCall, DXCCEntity(getDXCCEntity(qso, setup, prefixes)))
}

func sarlArea(call callsign.Callsign, dxccEntity DXCCEntity) string {
	// according to the SARL HF Phone, Digital and CW Contests §6.4, see https://mysarl.org.za/contest-resources/

	// the prefix database resolves ZS7 to Antarctica, not to South Africa
	if strings.HasPrefix(strings.ToUpper(call.BaseCall), "ZS7") {
		return "8"
	}

	switch strings.ToLower(string(dxccEntity)) {
	case "zs":
		area := callAreaDigit(call)
		if area >= "1" && area <= "6" {
			return area
		}
		return "9"
	case "v5":
		return "7"
	case "3da", "7p", "7q", "9j", "c9", "a2", "d2", "z2", "zd7", "zd9", "zs8", "fr", "3b8", "5r", "fh", "d6":
		return "8"
	default:
		return "9"
	}
}

func getBalkanPrefix(qso QSO, _ Setup, _ PrefixDatabase) string {
	return balkanPrefix(qso.TheirCall)
}

var balkanCallsignPrefixes = []string{
	"ZC4",
	"4O", "5B", "9A", "C4", "E7", "ER", "H2", "J4", "LZ", "P3", "S5", "SV", "SW", "SX",
	"SY", "SZ", "TA", "TB", "TC", "YM", "YO", "YP", "YQ", "YR", "YT", "YU", "Z3", "Z6", "ZA",
}

func balkanPrefix(call callsign.Callsign) string {
	// according to the Balkan HF Contest rules §2 and §11, see https://bfra.bg/en/node/464
	base := strings.ToUpper(call.BaseCall)
	if len(base) < 3 {
		return ""
	}

	isBalkan := false
	for _, prefix := range balkanCallsignPrefixes {
		if strings.HasPrefix(base, prefix) {
			isBalkan = true
			break
		}
	}
	if !isBalkan {
		return ""
	}

	// operation from another call area counts for that area
	suffix := strings.ToUpper(call.Suffix)
	if len(suffix) == 1 && suffix[0] >= '0' && suffix[0] <= '9' {
		return base[:2] + suffix
	}
	return base[:3]
}
