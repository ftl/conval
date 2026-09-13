/*
This file contains the implementation of the contest specific properties:
- EUDXCC (https://eudxcc.altervista.org/)
- VERON, the Dutch amateur radio association (https://www.veron.nl)
- SAC, the Scandinavian Activity Contest (https://www.sactest.net)
- SARL, the South African Radio League (https://mysarl.org.za)
- BFRA, the Bulgarian Federation of Radio Amateurs (https://bfra.bg)
- RSGB, the Radio Society of Great Britain (https://www.rsgbcc.org)
- the WTZC Committee (https://wtzc-contest.com)
*/
package conval

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/ftl/hamradio/callsign"
)

func init() {
	commonPropertyValidators[EURegionProperty] = RegexpValidator(validEURegion, "EU region")
	commonPropertyValidators[PAProvinceProperty] = RegexpValidator(validPAProvince, "PA province")
	commonPropertyValidators[CommonwealthHQProperty] = RegexpValidator(validCommonwealthHQ, "Commonwealth HQ")

	commonPropertyGetters[EURegionProperty] = getTheirExchangeProperty(EURegionProperty)
	commonPropertyGetters[PAProvinceProperty] = getTheirExchangeProperty(PAProvinceProperty)
	commonPropertyGetters[VeronEntityProperty] = PropertyGetterFunc(getVeronEntity)
	commonPropertyGetters[SACAreaProperty] = PropertyGetterFunc(getSACArea)
	commonPropertyGetters[SARLAreaProperty] = PropertyGetterFunc(getSARLArea)
	commonPropertyGetters[BalkanPrefixProperty] = PropertyGetterFunc(getBalkanPrefix)
	commonPropertyGetters[CommonwealthHQProperty] = getTheirExchangeProperty(CommonwealthHQProperty)
	commonPropertyGetters[CommonwealthAreaProperty] = PropertyGetterFunc(getCommonwealthArea)
	commonPropertyGetters[WTZCOffsetHoursProperty] = PropertyGetterFunc(getWTZCOffsetHours)

	myPropertyGetters[PAProvinceProperty] = getMyExchangeProperty(PAProvinceProperty)
	myPropertyGetters[VeronEntityProperty] = PropertyGetterFunc(getMyVeronEntity)
	myPropertyGetters[CommonwealthAreaProperty] = PropertyGetterFunc(getMyCommonwealthArea)
}

const (
	EURegionProperty         Property = "eu_region"
	PAProvinceProperty       Property = "pa_province"
	VeronEntityProperty      Property = "veron_entity"
	SACAreaProperty          Property = "sac_area"
	SARLAreaProperty         Property = "sarl_area"
	BalkanPrefixProperty     Property = "balkan_prefix"
	CommonwealthHQProperty   Property = "commonwealth_hq"
	CommonwealthAreaProperty Property = "commonwealth_area"
	WTZCOffsetProperty       Property = "wtzc_offset"
	WTZCOffsetHoursProperty  Property = "wtzc_offset_hours"
)

var (
	validEURegion       = regexp.MustCompile(`AT0[1-9]|BE[01][0-9]|BG0[1-6]|CZ[01][0-9]|CY0[1-5]|DK0[1-6]|EE0[1-5]|FI[01][0-9]|FR[0-2][0-9]|DE[01][0-9]|GR[01][0-9]|HU0[1-7]|IE0[1-4]|IT[0-2][0-9]|LV0[1-6]|LT0[1-5]|LX01|MT0[1-5]|NL[01][0-9]|PL[01][0-9]|RO0[1-8]|SK0[1-8]|SI0[1-6]|ES[01][1-9]|SE[0-2][1-9]`)
	validPAProvince     = regexp.MustCompile(`DR|FL|FR|GD|GR|LB|NB|NH|OV|UT|ZH|ZL`)
	validCommonwealthHQ = regexp.MustCompile(`HQ`)
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

func getMyCommonwealthArea(_ QSO, setup Setup, _ PrefixDatabase) string {
	return commonwealthArea(setup.MyCall, setup.MyCountry)
}

func getCommonwealthArea(qso QSO, setup Setup, prefixes PrefixDatabase) string {
	if qso.TheirExchange[CommonwealthHQProperty] != "" {
		// each HQ station counts as an additional call area
		return strings.ToUpper(qso.TheirCall.String())
	}
	return commonwealthArea(qso.TheirCall, DXCCEntity(getDXCCEntity(qso, setup, prefixes)))
}

// the entities of the Commonwealth call area list, see https://www.rsgbcc.org/hf/information/codes.shtml
var commonwealthEntities = []string{
	"1s", "3b6", "3b8", "3b9", "3d2", "3da", "4s", "5b", "5h", "5n", "5v", "5w", "5x", "5z", "6y", "7p", "7q",
	"8p", "8q", "8r", "9g", "9h", "9j", "9l", "9m2", "9m6", "9v", "9x", "9y", "a2", "a3", "ap", "c2", "c5",
	"c6", "c9", "cy0", "cy9", "e5/n", "e5/s", "e6", "g", "gd", "gi", "gj", "gm", "gu", "gw", "h4", "h40", "j3",
	"j6", "j7", "j8", "p2", "s2", "s7", "t2", "t30", "t31", "t32", "t33", "tj", "tr", "v2", "v3", "v4", "v5",
	"v8", "ve", "vk", "vk0h", "vk9c", "vk9m", "vk9n", "vk9w", "vk9x", "vp2e", "vp2m", "vp2v", "vp5", "vp6",
	"vp6/d", "vp8", "vp9", "vq9", "vu", "vu4", "vu7", "yj", "zb", "zc4", "zd7", "zd8", "zd9", "zf", "zk3", "zl",
	"zl7", "zl8", "zl9", "zs", "zs8",
}

func commonwealthArea(call callsign.Callsign, dxccEntity DXCCEntity) string {
	entity := strings.ToLower(string(dxccEntity))
	if !slices.Contains(commonwealthEntities, entity) {
		return ""
	}

	// the call areas of VE, VK, ZL and ZS are split in the same way as for the VERON contests
	return VeronEntity(call, dxccEntity)
}

func getWTZCOffsetHours(qso QSO, setup Setup, _ PrefixDatabase) string {
	// the shortest distance around the 24 hour clock in full hours, see https://wtzc-contest.com/rules section 5.1
	myOffset, myOK := utcOffset(qso.MyExchange[WTZCOffsetProperty])
	if !myOK {
		myOffset, myOK = utcOffset(setup.MyExchange[WTZCOffsetProperty])
	}
	theirOffset, theirOK := utcOffset(qso.TheirExchange[WTZCOffsetProperty])
	if !myOK || !theirOK {
		return ""
	}

	distance := myOffset - theirOffset
	if distance < 0 {
		distance = -distance
	}
	if distance > 12*60 {
		distance = 24*60 - distance
	}
	return strconv.Itoa(distance / 60)
}

var utcOffsetExpression = regexp.MustCompile(`^([0-9]{2})([0-9]{2})([EWZ])$`)

func utcOffset(value string) (int, bool) {
	matches := utcOffsetExpression.FindStringSubmatch(strings.ToUpper(strings.TrimSpace(value)))
	if matches == nil {
		return 0, false
	}
	hours, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, false
	}
	minutes, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, false
	}

	offset := hours*60 + minutes
	if matches[3] == "W" {
		offset = -offset
	}
	return offset, true
}
