package model

import (
	. "github.com/go-yaaf/yaaf-common/entity"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/enums"
)

// User represents a human / system operator that has access to the system, and can perform operations
// User authentication is done by an external identity provider
// @Entity: user
type User struct {
	BaseEntityEx
	Name       string         `json:"name"`       // User name
	Email      string         `json:"email"`      // User email
	Mobile     string         `json:"mobile"`     // User mobile phone number (for notification and validation)
	Type       UserTypeCode   `json:"type"`       // User type: UNDEFINED | SYSADMIN | SUPPORT | USER
	Roles      UserRoleFlag   `json:"roles"`      // User roles flags
	Groups     []string       `json:"groups"`     // User permissions groups
	Status     UserStatusCode `json:"status"`     // User status: UNDEFINED | PENDING | ACTIVE |  BLOCKED | SUSPENDED
	LastSignIn Timestamp      `json:"lastSignIn"` // User last successful sign in timestamp [epoch time milliseconds]
}

func (u *User) TABLE() string { return "user" }
func (u *User) NAME() string {
	if len(u.Name) > 0 {
		return u.Name
	} else {
		return u.Email
	}
}

// NewUser is a factory method to create new instance
func NewUser() Entity {
	return &User{BaseEntityEx: BaseEntityEx{CreatedOn: Now(), UpdatedOn: Now(), Props: make(Json)}, Groups: make([]string, 0)}
}

func (u *User) GetRoles() []UserRoleFlag {
	return SplitRoles(u.Roles)
}
