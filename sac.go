/*
This file contains the implementation of the specific things for the Scandinavian Activity Contest (https://www.sactest.net).
*/
package conval

import (
	"regexp"
	"strings"

	"github.com/ftl/hamradio/callsign"
)

func init() {
	commonPropertyGetters[SACAreaProperty] = PropertyGetterFunc(getSACArea)
}

const (
	SACAreaProperty Property = "sac_area"
)

func getSACArea(qso QSO, setup Setup, prefixes PrefixDatabase) string {
	return sacArea(qso.TheirCall, DXCCEntity(getDXCCEntity(qso, setup, prefixes)))
}

func sacArea(call callsign.Callsign, dxccEntity DXCCEntity) string {
	entity := strings.ToLower(string(dxccEntity))
	if entity == "" {
		return ""
	}
	return entity + sacCallAreaNumber(call)
}

var sacCallAreaExpression = regexp.MustCompile("[0-9]")

// see https://www.sactest.net/blog/scandinavian-activity-contest-2025-rules/ §8.2
func sacCallAreaNumber(call callsign.Callsign) string {
	prefix := call.Prefix
	if prefix == "" {
		prefix = call.BaseCall
	}
	if len(prefix) > 1 {
		// the first character may be a digit that belongs to the prefix, not to the call area
		prefix = prefix[1:]
	}

	number := sacCallAreaExpression.FindString(prefix)
	if number == "" {
		return "0"
	}
	return number
}
