package model

// UserStatusCode represents the user status: PENDING | ACTIVE | BLOCKED ...
// @Enum
type UserStatusCode = int

// UserStatusCodes the list of user status values
// @EnumValuesFor: UserStatusCode
type userStatusCode struct {

	// Undefined [0]
	UNDEFINED UserStatusCode `value:"0"`

	// User is registered and pending verification [1]
	PENDING UserStatusCode `value:"1"`

	// Active user in the system [2]
	ACTIVE UserStatusCode `value:"2"`

	// Blocked user (only account system can unblock the user) [3]
	BLOCKED UserStatusCode `value:"3"`

	// Suspended user (about to be deleted) [4]
	SUSPENDED UserStatusCode `value:"4"`

	IsValid func(int) bool
	String  func(int) string
}

var UserStatusCodes = &userStatusCode{
	UNDEFINED: 0,
	PENDING:   1,
	ACTIVE:    2,
	BLOCKED:   3,
	SUSPENDED: 4,
	IsValid:   isValidUserStatusCode,
	String:    stringUserStatusCode,
}

func isValidUserStatusCode(code int) bool {
	return code >= 0 && code <= 4
}

var userStatusCodes = []string{
	"UNDEFINED",
	"PENDING",
	"ACTIVE",
	"BLOCKED",
	"SUSPENDED",
}

func stringUserStatusCode(code int) string {
	if isValidUserStatusCode(code) {
		return userStatusCodes[code]
	} else {
		return "UNKNOWN"
	}
}
