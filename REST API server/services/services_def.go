package services

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	. "github.com/go-yaaf/yaaf-common/entity"
	"github.com/go-yaaf/yaaf-common/logger"
	"github.com/go-yaaf/yaaf-examples/rest-api/common"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/common"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/entities"
)

const (
	actionCreate = "Create"
	actionUpdate = "Update"
	actionDelete = "Delete"
)

type BaseService struct {
	ServiceName string
}

// Return formatted service error
func (s *BaseService) serviceError(method string, err error) error {
	if err == nil {
		return nil
	} else {
		logger.Error("[%s:%s]: %s", s.ServiceName, method, err.Error())
		return err //fmt.Errorf("%s:%s error: %s", s.ServiceName, method, err.Error())
	}
}

func (s *BaseService) serviceErrorEx(method string, code int, errText string) Error {
	logger.Error("[%s:%s]: %s", s.ServiceName, method, errText)
	return NewError(code, fmt.Sprintf("%s:%s error: %s", s.ServiceName, method, errText))
}

// Return custom formatted service error
func (s *BaseService) serviceErrorf(method string, message string, args ...any) error {
	errMsg := fmt.Sprintf(message, args...)
	logger.Error("[%s:%s]: %s", s.ServiceName, method, errMsg)
	return fmt.Errorf("%s:%s error: %s", s.ServiceName, method, errMsg)
}

// Calculate number of pages in the query based on total items and page size
func (s *BaseService) calcPages(total int64, size int) int {
	last := 0
	if total%int64(size) > 0 {
		last = 1
	}
	rem := float64(total / int64(size))
	return int(math.Round(rem)) + last
}

// Save user action to audit log
func (s *BaseService) auditLog(td *TokenData, entity Entity, action string, before, after any) {

	if td == nil || entity == nil {
		return
	}

	log := NewAuditLog()
	log.(*AuditLog).Id = IDN()
	log.(*AuditLog).UserId = td.SubjectId
	log.(*AuditLog).UserType = td.SubjectType
	log.(*AuditLog).Action = action
	log.(*AuditLog).ItemType = entity.TABLE()
	log.(*AuditLog).ItemId = entity.ID()
	log.(*AuditLog).ItemName = entity.NAME()

	log.(*AuditLog).BeforeChange = s.serializeChanges(before)
	log.(*AuditLog).AfterChange = s.serializeChanges(after)

	// Save log entry to the database
	_, _ = common.GetServiceHub().Database.Insert(log)
}

func (s *BaseService) serializeChanges(changes interface{}) (changesJson string) {
	changesJson = "{}"
	if bytes, err := json.Marshal(s.getStruct(changes)); err == nil {
		if bytes != nil {
			changesJson = string(bytes)
		}
	}
	return
}

func (s *BaseService) getStruct(object any) (result interface{}) {
	if object == nil {
		return nil
	}
	switch v := object.(type) {
	case Entity:
		return v
	default:
		return struct {
			Value any `json:"value"`
		}{v}
	}
}

// Get unique identifiers from list
func (s *BaseService) getUniqueIds(entities []Entity) []string {
	strMap := make(map[string]string)
	for _, ent := range entities {
		strMap[ent.ID()] = ent.ID()
	}
	vsm := make([]string, 0)
	for k, _ := range strMap {
		vsm = append(vsm, k)
	}
	return vsm
}

// Get map of id->entity from entity list
func (s *BaseService) getEntityMap(entities []Entity) map[string]Entity {
	result := make(map[string]Entity)
	for _, ent := range entities {
		result[ent.ID()] = ent
	}
	return result
}

func ToAnyVariadic[T any](items []T) (result []any) {
	for _, t := range items {
		result = append(result, t)
	}
	return result
}

// Remove all spaces and dashes from phone number
func (s *BaseService) stripPhone(phone string) string {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	return phone
}
