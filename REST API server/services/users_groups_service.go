package services

import (
	"fmt"
	"sync"

	. "github.com/go-yaaf/yaaf-common/database"
	. "github.com/go-yaaf/yaaf-common/entity"

	. "github.com/go-yaaf/yaaf-examples/rest-api/common"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/common"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/entities"
	. "github.com/go-yaaf/yaaf-examples/rest-api/utils"
)

var groupsServiceOnce sync.Once
var groupsServiceInst *GroupsService = nil

type GroupsService struct {
	BaseService
	sh *ServiceHub // Service hub
}

// GetGroupsService factory function
func GetGroupsService(sh *ServiceHub) *GroupsService {
	groupsServiceOnce.Do(func() {
		if groupsServiceInst == nil {
			groupsServiceInst = &GroupsService{BaseService: BaseService{ServiceName: "GroupsService"}, sh: sh}
		}
	})
	return groupsServiceInst
}

// Create new group in the system
func (s *GroupsService) Create(td *TokenData, entity Entity) (Entity, error) {

	ent := entity.(*UsersGroup)

	// Override system fields,
	if len(ent.Id) == 0 {
		ent.Id = TokenUtils().NanoID()
	}
	ent.CreatedOn = Now()
	ent.UpdatedOn = Now()
	ent.Props = nil

	ent.Members = nil

	if updated, er := s.sh.Database.Insert(ent); er != nil {
		return nil, s.serviceError("Create", er)
	} else {
		s.auditLog(td, ent, actionCreate, nil, updated)
		return updated, nil
	}
}

// Update existing group in the system
func (s *GroupsService) Update(td *TokenData, entity Entity) (Entity, error) {

	ent := entity.(*UsersGroup)

	// Get existing group
	existing, err := s.sh.Database.Get(NewUsersGroup, ent.Id)
	if err != nil {
		return nil, s.serviceError("Update", err)
	}

	// Override system fields,
	ent.CreatedOn = existing.(*UsersGroup).CreatedOn
	ent.UpdatedOn = Now()
	ent.Props = nil

	if updated, er := s.sh.Database.Update(ent); er != nil {
		return nil, s.serviceError("Update", er)
	} else {
		s.auditLog(td, ent, actionUpdate, existing, updated)
		return updated, nil
	}
}

// Delete group
func (s *GroupsService) Delete(td *TokenData, id string) (err error) {

	// Get existing group
	var existing Entity
	if existing, err = s.sh.Database.Get(NewUsersGroup, id); err != nil {
		return s.serviceError("Delete", err)
	}

	if err = s.sh.Database.Delete(NewUsersGroup, id); err != nil {
		return s.serviceError("Delete", err)
	} else {
		s.auditLog(td, existing, actionDelete, existing, nil)
		return nil
	}
}

// Get single group by id
func (s *GroupsService) Get(td *TokenData, id string) (Entity, error) {
	if ent, err := s.sh.Database.Get(NewUsersGroup, id); err != nil {
		return nil, fmt.Errorf("[%s]::Get: %v", s.ServiceName, err)
	} else {
		return ent, nil
	}
}

// GroupsFindParams Query params aggregator for find commands service
type GroupsFindParams struct {
	Search string // Filter by free text search (using * wildcard)
	Sort   string // Sort descriptor (field name with suffix +/- for sort order)
	Page   int    // Page number for pagination
	Size   int    // Page size: number of items per page
}

// Find list of groups by filter
func (s *GroupsService) Find(p GroupsFindParams) (entities []Entity, total int64, pages int, error error) {
	if entities, total, error = s.sh.Database.Query(NewUsersGroup).
		MatchAny(
			F("id").Like(p.Search),
			F("name").Like(p.Search),
		).
		Page(p.Page).
		Limit(p.Size).
		Sort(p.Sort).
		Find(); error == nil {
		pages = s.calcPages(total, p.Size)
	} else {
		error = s.serviceError("Find", error)
	}
	return
}
