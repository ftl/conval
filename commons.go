package conval

// Common Exchange Validators

func ValidateRST(exchange string) error {
	return nil // TODO implement
}

func ValidateSerialNumber(exchange string) error {
	return nil // TODO implement
}

func ValidateMemberNumber(exchange string) error {
	return nil // TODO implement
}

func ValidateNoMember(exchange string) error {
	return nil // TODO implement
}

func ValidateCQZone(exchange string) error {
	return nil // TODO implement
}

func ValidateITUZone(exchange string) error {
	return nil // TODO implement
}

// Common Property Getters

func GetCQZone(qso QSO) string {
	exchange, ok := qso.TheirExchange[CQZoneProperty]
	if ok {
		return exchange
	}
	// TODO get CQ zone from database
	return ""
}

func GetITUZone(qso QSO) string {
	exchange, ok := qso.TheirExchange[ITUZoneProperty]
	if ok {
		return exchange
	}
	// TODO get ITU zone from database
	return ""
}

func GetDXCCEntity(qso QSO) string {
	return "" // TODO implement
}

func GetWPXPrefix(qso QSO) string {
	return "" // TODO implement
}

func GetTheirExchangeProperty(property Property) PropertyGetter {
	return PropertyGetterFunc(func(qso QSO) string {
		return qso.TheirExchange[property]
	})
}
