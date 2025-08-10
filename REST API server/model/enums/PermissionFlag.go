package model

// PermissionFlag represents combination of permissions: READ | CREATE | UPDATE | DELETE | MANAGE
// @Enum
type PermissionFlag = int

// PermissionFlag represents the list of permissions
// @EnumValuesFor: PermissionFlag
type permissionFlag struct {

	// Undefined [0]
	UNDEFINED PermissionFlag `value:"0"`

	// Read [1]
	READ PermissionFlag `value:"1"`

	// Create [2]
	CREATE PermissionFlag `value:"2"`

	// Update [4]
	UPDATE PermissionFlag `value:"4"`

	// Delete [8]
	DELETE PermissionFlag `value:"8"`

	// Manage [16]
	MANAGE PermissionFlag `value:"16"`

	// All permissions combined
	ALL PermissionFlag `value:"31"`

	IsValid func(int) bool
	String  func(int) string
}

var PermissionFlags = &permissionFlag{
	UNDEFINED: 0,  // Undefined [0]
	READ:      1,  // Read [1]
	CREATE:    2,  // Create  [2]
	UPDATE:    4,  // Update [4]
	DELETE:    8,  // Delete [8]
	MANAGE:    16, // Manage [16]
	ALL:       31, // All permissions [2047]
	IsValid:   isValidPermissionFlag,
	String:    stringPermissionFlag,
}

func isValidPermissionFlag(code int) bool {
	return true
}

var permissionFlags = []string{
	"UNDEFINED",
	"READ",
	"CREATE",
	"UPDATE",
	"DELETE",
	"MANAGE",
	"ALL",
}

func stringPermissionFlag(code int) string {
	if isValidPermissionFlag(code) {
		return permissionFlags[code]
	} else {
		return "UNKNOWN"
	}
}

func SplitPermissions(code int) []PermissionFlag {
	permissions := make([]PermissionFlag, 0)

	if code&PermissionFlags.CREATE == PermissionFlags.CREATE {
		permissions = append(permissions, PermissionFlags.CREATE)
	}
	if code&PermissionFlags.UPDATE == PermissionFlags.UPDATE {
		permissions = append(permissions, PermissionFlags.UPDATE)
	}
	if code&PermissionFlags.READ == PermissionFlags.READ {
		permissions = append(permissions, PermissionFlags.READ)
	}
	if code&PermissionFlags.DELETE == PermissionFlags.DELETE {
		permissions = append(permissions, PermissionFlags.DELETE)
	}
	if code&PermissionFlags.MANAGE == PermissionFlags.MANAGE {
		permissions = append(permissions, PermissionFlags.MANAGE)
	}
	return permissions
}

func CombinePermissions(permissions []PermissionFlag) int {
	result := 0
	for _, permFlag := range permissions {
		result = result | permFlag
	}
	return result
}
