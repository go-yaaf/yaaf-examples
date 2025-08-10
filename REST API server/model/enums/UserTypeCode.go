package model

// UserTypeCode represents the user type: SYSADMIN | SUPPORT | USER ...
// @Enum
type UserTypeCode = int

// UserTypeCodes represents the list of user types
// @EnumValuesFor: UserTypeCode
type userTypeCode struct {

	// Undefined [0]
	UNDEFINED UserTypeCode `value:"0"`

	// System administrator has access to all accounts and permissions to perform all actions [1]
	SYSADMIN UserTypeCode `value:"1"`

	// Support user has view permissions only for all accounts that enabled option Enable Support [2]
	SUPPORT UserTypeCode `value:"2"`

	// Account user - has access to specific accounts with role based access control [3]
	USER UserTypeCode `value:"3"`

	// Service Account - to be used by other systems to perform actions using the API (can't login as a user to the portal) [4]
	SERVICE UserTypeCode `value:"4"`

	IsValid func(int) bool
	String  func(int) string
}

var UserTypeCodes = &userTypeCode{
	UNDEFINED: 0, // Undefined [0]
	SYSADMIN:  1, // System administrator has access to all accounts and permissions to perform all actions [1]
	SUPPORT:   2, // Support user has view permissions only for all accounts that enabled option Enable Support [2]
	USER:      3, // Account user - has access to specific accounts with role based access control [3]
	SERVICE:   4, // Account service - to be used by other systems to perform actions using the API (can't login as a user to the portal) [4]
	IsValid:   isValidUserTypeCode,
	String:    stringUserTypeCode,
}

func isValidUserTypeCode(code int) bool {
	return code >= 0 && code <= 5
}

var userTypeCodes = []string{
	"UNDEFINED",
	"SYSADMIN",
	"SUPPORT",
	"USER",
	"SERVICE",
}

func stringUserTypeCode(code int) string {
	if isValidUserTypeCode(code) {
		return userTypeCodes[code]
	} else {
		return "UNKNOWN"
	}
}
