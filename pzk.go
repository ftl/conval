/*
This file contains the implementation of the specific things for contests announced by the Polish Amateur Radio Union PZK.
*/
package conval

import (
	"regexp"
)

func init() {
	PropertyValidators[PZKProvinceProperty] = RegexpValidator(validPZKProvince, "PZK province")

	PropertyGetters[PZKProvinceProperty] = getTheirExchangeProperty(PZKProvinceProperty)
}

const (
	PZKProvinceProperty Property = "pzk_province"
)

var (
	validPZKProvince = regexp.MustCompile(`B|C|D|F|G|J|K|L|M|O|P|R|S|U|W|Z`)
)
