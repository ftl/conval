/*
This file contains the implementation of the specific things for contests announced by the CQ magazine.
*/
package conval

import (
	"regexp"
	"unicode"

	"github.com/ftl/hamradio/callsign"
)

func init() {
	PropertyGetters[WPXPrefixProperty] = PropertyGetterFunc(getWPXPrefix)
}

const (
	WPXPrefixProperty Property = "wpx_prefix"
)

func getWPXPrefix(qso QSO) string {
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
