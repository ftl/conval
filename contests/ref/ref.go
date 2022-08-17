/*
The package ref implements the specific things for contests announced by the REF.
*/
package ref

import "github.com/ftl/conval"

func init() {
	conval.ExchangeValidators[DepartmentExchange] = conval.ExchangeValidatorFunc(ValidateDepartment)
	conval.PropertyGetters[DepartmentProperty] = conval.PropertyGetterFunc(GetDepartment)
}

const (
	DepartmentExchange conval.Exchange = "ref_department"

	DepartmentProperty conval.Property = "ref_department"
)

func ValidateDepartment(exchange string) error {
	return nil // TODO implement
}

func GetDepartment(qso conval.QSO) string {
	return "" // TODO implement
}
