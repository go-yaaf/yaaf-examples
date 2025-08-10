package model

// AccountStatusCode represents the account status: ACTIVE | INACTIVE | BLOCKED ...
// @Enum
type AccountStatusCode = int

// List of account status values
// @EnumValuesFor: AccountStatusCode
type accountStatusCode struct {
	// Undefined [0]
	UNDEFINED AccountStatusCode `value:"0"`

	// Active account in the system [1]
	ACTIVE AccountStatusCode `value:"1"`

	// Inactive account in the system [2]
	INACTIVE AccountStatusCode `value:"2"`

	// Blocked account [3]
	BLOCKED AccountStatusCode `value:"3"`

	// Suspended account (about to be deleted) [4]
	SUSPENDED AccountStatusCode `value:"4"`

	IsValid func(int) bool
	String  func(int) string
}

var AccountStatusCodes = &accountStatusCode{
	UNDEFINED: 0, // Undefined [0]
	ACTIVE:    1, // Active account in the system [1]
	INACTIVE:  2, // Inactive account in the system [2]
	BLOCKED:   3, // Blocked account [3]
	SUSPENDED: 4, // Suspended account (about to be deleted) [4]
	IsValid:   isValidAccountStatusCode,
	String:    stringAccountStatusCode,
}

func isValidAccountStatusCode(code int) bool {
	return code >= 0 && code <= 4
}

var accountStatusCodes = []string{
	"UNDEFINED",
	"ACTIVE",
	"INACTIVE",
	"BLOCKED",
	"SUSPENDED",
}

func stringAccountStatusCode(code int) string {
	if isValidAccountStatusCode(code) {
		return accountStatusCodes[code]
	} else {
		return "UNKNOWN"
	}
}
