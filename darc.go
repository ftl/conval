/*
This file contains the implementation of the specific things for contests announced by the DARC.
*/
package conval

func init() {
	PropertyValidators[WAGDOKProperty] = PropertyValidatorFunc(validateWAGDOK)

	PropertyGetters[WAEEntityProperty] = PropertyGetterFunc(getWAEEntity)
	PropertyGetters[WAGDistrictProperty] = PropertyGetterFunc(getWAGDistrict)
}

const (
	WAGDOKProperty      Property = "wag_dok"
	WAEEntityProperty   Property = "wae_property"
	WAGDistrictProperty Property = "wag_district"
)

func validateWAGDOK(exchange string) error {
	return nil // TODO implement
}

func getWAEEntity(qso QSO) string {
	return "" // TODO implement
}

func getWAGDistrict(qso QSO) string {
	return "" // TODO implement
}
