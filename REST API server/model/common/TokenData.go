package model

import (
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/enums"
)

// TokenData model represents user in account which is encrypted with the JWT token
// @Data
type TokenData struct {
	SubjectId   string         `json:"subjectId"`   // Authenticated subject ID (can be user, or service account)
	SubjectType UserTypeCode   `json:"subjectType"` // Subject type: UNDEFINED | SYSADMIN | USER | SERVICE_ACCOUNT
	Status      UserStatusCode `json:"status"`      // User status: UNDEFINED | PENDING | ACTIVE | BLOCKED | SUSPENDED
	ExpiresIn   int64          `json:"expiresIn"`   // Token expiration [Epoch milliseconds Timestamp]
}
