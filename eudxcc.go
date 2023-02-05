/*
This file contains the implementation of the specific things for contests announced by the EUDXCC (https://eudxcc.altervista.org/).
*/
package conval

import (
	"regexp"
)

func init() {
	PropertyValidators[EURegionProperty] = RegexpValidator(validEURegion, "EU region")

	PropertyGetters[EURegionProperty] = getTheirExchangeProperty(EURegionProperty)
}

const (
	EURegionProperty Property = "eu_region"
)

var (
	validEURegion = regexp.MustCompile(`AT0[1-9]|BE[01][0-9]|BG0[1-6]|CZ[01][0-9]|CY0[1-5]|DK0[1-6]|EE0[1-5]|FI[01][0-9]|FR[0-2][0-9]|DE[01][0-9]|GR[01][0-9]|HU0[1-7]|IE0[1-4]|IT[0-2][0-9]|LV0[1-6]|LT0[1-5]|LX01|MT0[1-5]|NL[01][0-9]|PL[01][0-9]|RO0[1-8]|SK0[1-8]|SI0[1-6]|ES[01][1-9]|SE[0-2][1-9]`)
)
