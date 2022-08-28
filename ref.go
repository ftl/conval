/*
This file contains the implementation of the specific things for contests announced by the REF.
*/
package conval

func init() {
	PropertyValidators[REFDepartmentProperty] = PropertyValidatorFunc(validateREFDepartment)

	PropertyGetters[REFDepartmentProperty] = PropertyGetterFunc(getREFDepartment)
}

const (
	REFDepartmentProperty Property = "ref_department"
)

func validateREFDepartment(exchange string) error {
	return nil // TODO implement
}

func getREFDepartment(qso QSO) string {
	return "" // TODO implement
}
