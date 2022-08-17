/*
The package darc implements the specific things for contests announced by the DARC.
*/
package darc

import "github.com/ftl/conval"

func init() {
	conval.ExchangeValidators[WAGDOKExchange] = conval.ExchangeValidatorFunc(ValidateWAGDOK)
	conval.PropertyGetters[WAEEntityProperty] = conval.PropertyGetterFunc(GetWAEEntity)
	conval.PropertyGetters[WAGDistrictProperty] = conval.PropertyGetterFunc(GetWAGDistrict)
}

const (
	WAGDOKExchange conval.Exchange = "wag_dok"

	WAEEntityProperty   conval.Property = "wae_property"
	WAGDistrictProperty conval.Property = "wag_district"
)

func ValidateWAGDOK(exchange string) error {
	return nil // TODO implement
}

func GetWAEEntity(qso conval.QSO) string {
	return "" // TODO implement
}

func GetWAGDistrict(qso conval.QSO) string {
	return "" // TODO implement
}
