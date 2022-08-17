package conval

// Common Exchange Validators

func ValidateRST(exchange string) error {
	return nil // TODO implement
}

func ValidateSerial(exchange string) error {
	return nil // TODO implement
}

func ValidateMemberNumber(exchange string) error {
	return nil // TODO implement
}

func ValidateCQZone(exchange string) error {
	return nil // TODO implement
}

func ValidateITUZone(exchange string) error {
	return nil // TODO implement
}

func ValidateNoMember(exchange string) error {
	return nil // TODO implement
}

// Common Property Getters

func GetCQZone(qso QSO) string {
	return "" // TODO implement
}

func GetITUZone(qso QSO) string {
	return "" // TODO implement
}

func GetDXCCEntity(qso QSO) string {
	return "" // TODO implement
}

func GetWPXPrefix(qso QSO) string {
	return "" // TODO implement
}
