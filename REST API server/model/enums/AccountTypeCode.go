package model

// AccountTypeCode represents the account type: DEMO | TRIAL | PARTNER | BUSINESS ...
// @Enum
type AccountTypeCode = int

// List of account status values
// @EnumValuesFor: AccountTypeCode
type accountTypeCode struct {
	// Undefined [0]
	UNDEFINED AccountTypeCode `value:"0"`

	// Demo account [1]
	DEMO AccountTypeCode `value:"1"`

	// Trial account [2]
	TRIAL AccountTypeCode `value:"2"`

	// Partner account [3]
	PARTNER AccountTypeCode `value:"3"`

	// Business account [4]
	BUSINESS AccountTypeCode `value:"4"`

	IsValid func(int) bool
	String  func(int) string
}

var AccountTypeCodes = &accountTypeCode{
	UNDEFINED: 0, // Undefined [0]
	DEMO:      1, // Demo account [1]
	TRIAL:     2, // Trial account [2]
	PARTNER:   3, // Partner account [3]
	BUSINESS:  4, // Business account [4]
	IsValid:   isValidAccountTypeCode,
	String:    stringAccountTypeCode,
}

func isValidAccountTypeCode(code int) bool {
	return code >= 0 && code <= 4
}

var accountTypeCodes = []string{
	"UNDEFINED",
	"DEMO",
	"TRIAL",
	"PARTNER",
	"BUSINESS",
}

func stringAccountTypeCode(code int) string {
	if isValidAccountTypeCode(code) {
		return accountTypeCodes[code]
	} else {
		return "UNKNOWN"
	}
}
