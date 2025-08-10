package model

import (
	. "github.com/go-yaaf/yaaf-common/entity"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/common"
)

// Contact entity is a billing account in the system
// @Entity: contact
type Contact struct {
	BaseEntityEx
	AccountId   string   `json:"accountId"`   // Related billing account ID
	Name        string   `json:"name"`        // Contact name
	Description string   `json:"description"` // Contact description
	Mobile      string   `json:"mobile"`      // Mobile phone
	Email       string   `json:"email"`       // Email address
	Address     Address  `json:"address"`     // Contact address
	Groups      []string `json:"groups"`      // Contact groups
}

func (a *Contact) TABLE() string { return "contact" }
func (a *Contact) NAME() string  { return a.Name }
func (a *Contact) KEY() string   { return a.AccountId }

// NewContact is a factory method to create new instance
func NewContact() Entity {
	return &Contact{BaseEntityEx: BaseEntityEx{CreatedOn: Now(), UpdatedOn: Now(), Id: GUID(), Props: make(Json)}, Groups: make([]string, 0)}
}
