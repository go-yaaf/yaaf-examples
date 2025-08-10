package model

// UserRoleFlag represents combination of roles: STUDENT | PILOT | INSTRUCTOR | SALES | OPERATIONS | MANAGER
// @Enum
type UserRoleFlag = int

// UserRoleFlag represents the list of roles
// @EnumValuesFor: UserRoleFlag
type userRoleFlag struct {

	// Undefined [0]
	UNDEFINED UserRoleFlag `value:"0"`

	// Sales [8]
	SALES UserRoleFlag `value:"8"`

	// Operations [16]
	OPERATIONS UserRoleFlag `value:"16"`

	// Maintenance [32]
	MAINTENANCE UserRoleFlag `value:"32"`

	// Manager [1024]
	MANAGER UserRoleFlag `value:"1024"`

	// All roles combined
	ALL UserRoleFlag `value:"2047"`

	IsValid func(int) bool
	String  func(int) string
}

var UserRoleFlags = &userRoleFlag{
	UNDEFINED:   0,    // Undefined [0]
	SALES:       8,    // Sales [8]
	OPERATIONS:  16,   // Operations [16]
	MAINTENANCE: 32,   // Maintenance [32]
	MANAGER:     1024, // Manager [1024]
	ALL:         2047, // All roles [2047]
	IsValid:     isValidUserRoleFlag,
	String:      stringUserRoleFlag,
}

func isValidUserRoleFlag(code int) bool {
	return true
}

var userRoleFlags = []string{
	"UNDEFINED",
	"SALES",
	"OPERATIONS",
	"MAINTENANCE",
	"MANAGER",
	"ALL",
}

func stringUserRoleFlag(code int) string {
	return "NOT IMPLEMENTED"
}

func SplitRoles(code int) []UserRoleFlag {
	roles := make([]UserRoleFlag, 0)

	if code&UserRoleFlags.SALES == UserRoleFlags.SALES {
		roles = append(roles, UserRoleFlags.SALES)
	}
	if code&UserRoleFlags.OPERATIONS == UserRoleFlags.OPERATIONS {
		roles = append(roles, UserRoleFlags.OPERATIONS)
	}
	if code&UserRoleFlags.MAINTENANCE == UserRoleFlags.MAINTENANCE {
		roles = append(roles, UserRoleFlags.MAINTENANCE)
	}
	if code&UserRoleFlags.MANAGER == UserRoleFlags.MANAGER {
		roles = append(roles, UserRoleFlags.MANAGER)
	}
	return roles
}

func CombineRoles(roles []UserRoleFlag) int {
	role := 0
	for _, roleFlag := range roles {
		role = role | roleFlag
	}
	return role
}
