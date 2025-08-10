package model

import (
	. "github.com/go-yaaf/yaaf-common/entity"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/enums"
)

// AuditLog entity is a log entry in the audit log to track users / service account actions
// @Entity: audit_log
type AuditLog struct {
	BaseEntityEx
	UserId       string       `json:"userId"`       // User Id
	UserType     UserTypeCode `json:"userType"`     // User type: UNDEFINED | SYSADMIN | USER | SERVICE_ACCOUNT
	Action       string       `json:"action"`       // Action that was performed
	ItemType     string       `json:"itemType"`     // Item type
	ItemId       string       `json:"itemId"`       // Item Id
	ItemName     string       `json:"itemName"`     // Item Name
	BeforeChange string       `json:"beforeChange"` // Item value before change [Json]
	AfterChange  string       `json:"afterChange"`  // Item delta after change [Json]
}

func (a *AuditLog) TABLE() string { return "audit_log" }
func (a *AuditLog) NAME() string  { return a.ItemName }

// NewAuditLog is a factory method to create new instance
func NewAuditLog() Entity {
	return &AuditLog{BaseEntityEx: BaseEntityEx{CreatedOn: Now(), UpdatedOn: Now(), Id: GUID(), Props: make(Json)}}
}
