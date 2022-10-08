/*
This file contains the implementation of the common properties used for contest scoring.
*/
package conval

func init() {
	PropertyValidators[TheirRSTProperty] = PropertyValidatorFunc(validateRST)
	PropertyValidators[SerialNumberProperty] = PropertyValidatorFunc(validateSerialNumber)
	PropertyValidators[MemberNumberProperty] = PropertyValidatorFunc(validateMemberNumber)
	PropertyValidators[CQZoneProperty] = PropertyValidatorFunc(validateCQZone)
	PropertyValidators[ITUZoneProperty] = PropertyValidatorFunc(validateITUZone)
	PropertyValidators[NoMemberProperty] = PropertyValidatorFunc(validateNoMember)

	PropertyGetters[TheirRSTProperty] = getTheirExchangeProperty(TheirRSTProperty)
	PropertyGetters[SerialNumberProperty] = getTheirExchangeProperty(SerialNumberProperty)
	PropertyGetters[MemberNumberProperty] = getTheirExchangeProperty(MemberNumberProperty)
	PropertyGetters[NoMemberProperty] = getTheirExchangeProperty(NoMemberProperty)
	PropertyGetters[CQZoneProperty] = PropertyGetterFunc(getCQZone)
	PropertyGetters[ITUZoneProperty] = PropertyGetterFunc(getITUZone)
	PropertyGetters[DXCCEntityProperty] = PropertyGetterFunc(getDXCCEntity)
	PropertyGetters[WorkingConditionProperty] = PropertyGetterFunc(getCallsignWorkingCondition)
}

// Common Exchange Validators

func validateRST(exchange string) error {
	return nil // TODO implement
}

func validateSerialNumber(exchange string) error {
	return nil // TODO implement
}

func validateMemberNumber(exchange string) error {
	return nil // TODO implement
}

func validateNoMember(exchange string) error {
	return nil // TODO implement
}

func validateCQZone(exchange string) error {
	return nil // TODO implement
}

func validateITUZone(exchange string) error {
	return nil // TODO implement
}

// Common Property Getters

func getCQZone(qso QSO) string {
	exchange, ok := qso.TheirExchange[CQZoneProperty]
	if ok {
		return exchange
	}
	// TODO get CQ zone from database
	return ""
}

func getITUZone(qso QSO) string {
	exchange, ok := qso.TheirExchange[ITUZoneProperty]
	if ok {
		return exchange
	}
	// TODO get ITU zone from database
	return ""
}

func getDXCCEntity(qso QSO) string {
	if qso.TheirCountry != "" {
		return string(qso.TheirCountry)
	}
	// TODO get DXCC entity from database
	return ""
}

func getCallsignWorkingCondition(qso QSO) string {
	return qso.TheirCall.WorkingCondition
}

func getTheirExchangeProperty(property Property) PropertyGetter {
	return PropertyGetterFunc(func(qso QSO) string {
		return qso.TheirExchange[property]
	})
}
