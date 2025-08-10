package model

import (
	. "github.com/go-yaaf/yaaf-common/entity"
)

// UsersGroup represents a group of users to share permissions
// @Entity: users_group
type UsersGroup struct {
	BaseEntityEx
	Name    string   `json:"name"`    // Group name
	Email   string   `json:"email"`   // Group email
	Members []string `json:"members"` // List of group members (user Ids)
}

func (u *UsersGroup) TABLE() string { return "users_group" }
func (u *UsersGroup) NAME() string {
	if len(u.Name) > 0 {
		return u.Name
	} else {
		return u.Email
	}
}

// NewUsersGroup is a factory method to create new instance
func NewUsersGroup() Entity {
	return &UsersGroup{BaseEntityEx: BaseEntityEx{Id: GUID(), CreatedOn: Now(), UpdatedOn: Now(), Props: make(Json)}, Members: make([]string, 0)}
}
