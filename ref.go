/*
This file contains the implementation of the specific things for contests announced by the REF.
*/
package conval

import (
	"fmt"
	"strconv"
)

func init() {
	PropertyValidators[REFDepartmentProperty] = PropertyValidatorFunc(validateREFDepartment)

	PropertyGetters[REFDepartmentProperty] = getTheirExchangeProperty(REFDepartmentProperty)
}

const (
	REFDepartmentProperty Property = "ref_department"
)

func validateREFDepartment(exchange string) error {
	numeric, err := strconv.Atoi(exchange)
	if err == nil {
		if numeric < 0 || numeric > 95 {
			return fmt.Errorf("%s is not a valid numeric department", exchange)
		}
		return nil
	}
	switch exchange {
	case "2A", "2B", "FG", "FJ", "FH", "FK", "FM", "FO", "FP", "FR", "FT", "FW", "FY":
		return nil
	default:
		return fmt.Errorf("%s is not a valid department", exchange)
	}
}
