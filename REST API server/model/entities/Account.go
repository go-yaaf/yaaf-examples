package model

import (
	. "github.com/go-yaaf/yaaf-common/entity"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/enums"
)

// Account entity is a billing account in the system
// @Entity: account
type Account struct {
	BaseEntityEx
	Name        string            `json:"name"`        // Account name
	Description string            `json:"description"` // Account description
	Type        AccountTypeCode   `json:"type"`        // Account type:  STUDENT | PRIVATE | BUSINESS ...
	Status      AccountStatusCode `json:"status"`      // Account status: UNDEFINED | ACTIVE | INACTIVE | BLOCKED | SUSPENDED
	Phone       string            `json:"phone"`       // Office / Landline phone
	Mobile      string            `json:"mobile"`      // Mobile phone
	Email       string            `json:"email"`       // Email address
}

func (a *Account) TABLE() string { return "account" }
func (a *Account) NAME() string  { return a.Name }

// NewAccount is a factory method to create a new instance
func NewAccount() Entity {
	return &Account{BaseEntityEx: BaseEntityEx{CreatedOn: Now(), UpdatedOn: Now(), Id: GUID(), Props: make(Json)}}
}
