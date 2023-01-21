/*
This file contains the implementation of the specific things for contests announced by the Hungarian Radio Amateur Society MRASZ.
*/
package conval

import (
	"regexp"
)

func init() {
	PropertyValidators[HACountyProperty] = RegexpValidator(validHACounty, "HA county")

	PropertyGetters[HACountyProperty] = getTheirExchangeProperty(HACountyProperty)
}

const (
	HACountyProperty Property = "ha_county"
)

var (
	validHACounty = regexp.MustCompile(`BN|BA|BE|BO|CS|FE|GY|HB|HE|SZ|KO|NG|PE|SO|SA|TO|VA|VE|ZA|BP`)
)
