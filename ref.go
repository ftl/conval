/*
This file contains the implementation of the specific things for contests announced by the REF.
*/
package conval

func init() {
	PropertyValidators[REFDepartmentProperty] = PropertyValidatorFunc(validateREFDepartment)

	PropertyGetters[REFDepartmentProperty] = getTheirExchangeProperty(REFDepartmentProperty)
}

const (
	REFDepartmentProperty Property = "ref_department"
)

func validateREFDepartment(exchange string) error {
	return nil // TODO implement
}
