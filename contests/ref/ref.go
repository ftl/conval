/*
The package ref implements the specific things for contests announced by the REF.
*/
package ref

import "github.com/ftl/conval"

func init() {
	conval.PropertyValidators[DepartmentProperty] = conval.PropertyValidatorFunc(ValidateDepartment)
	conval.PropertyGetters[DepartmentProperty] = conval.PropertyGetterFunc(GetDepartment)
}

const (
	DepartmentProperty conval.Property = "ref_department"
)

func ValidateDepartment(exchange string) error {
	return nil // TODO implement
}

func GetDepartment(qso conval.QSO) string {
	return "" // TODO implement
}
