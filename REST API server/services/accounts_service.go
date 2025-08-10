package services

import (
	"fmt"
	"sync"

	. "github.com/go-yaaf/yaaf-common/database"
	. "github.com/go-yaaf/yaaf-common/entity"

	. "github.com/go-yaaf/yaaf-examples/rest-api/common"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/common"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/entities"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/enums"
	. "github.com/go-yaaf/yaaf-examples/rest-api/utils"
)

var accountsServiceOnce sync.Once
var accountsServiceInst *AccountsService = nil

type AccountsService struct {
	BaseService
	sh *ServiceHub // Service hub
}

// GetAccountsService factory function
func GetAccountsService(sh *ServiceHub) *AccountsService {
	accountsServiceOnce.Do(func() {
		if accountsServiceInst == nil {
			accountsServiceInst = &AccountsService{BaseService: BaseService{ServiceName: "AccountsService"}, sh: sh}
		}
	})
	return accountsServiceInst
}

// Create a new account in the system
func (s *AccountsService) Create(td *TokenData, entity Entity) (Entity, error) {

	ent := entity.(*Account)

	// Override system fields,
	ent.Id = TokenUtils().GUID()
	ent.CreatedOn = Now()
	ent.UpdatedOn = Now()
	ent.Props = nil

	// Strip phone numbers
	ent.Mobile = s.stripPhone(ent.Mobile)
	ent.Phone = s.stripPhone(ent.Phone)

	if updated, er := s.sh.Database.Insert(ent); er != nil {
		return nil, s.serviceError("Create", er)
	} else {
		s.auditLog(td, ent, actionCreate, nil, updated)
		return updated, nil
	}
}

// Update existing account in the system
func (s *AccountsService) Update(td *TokenData, entity Entity) (Entity, error) {

	ent := entity.(*Account)

	// Get existing account
	existing, err := s.sh.Database.Get(NewAccount, ent.Id)
	if err != nil {
		return nil, s.serviceError("Update", err)
	}

	// Override system fields,
	ent.CreatedOn = existing.(*Account).CreatedOn
	ent.UpdatedOn = Now()
	ent.Props = nil

	// Strip phone numbers
	ent.Mobile = s.stripPhone(ent.Mobile)
	ent.Phone = s.stripPhone(ent.Phone)

	if updated, er := s.sh.Database.Update(ent); er != nil {
		return nil, s.serviceError("Update", er)
	} else {
		s.auditLog(td, ent, actionUpdate, existing, updated)
		return updated, nil
	}
}

// Delete account
func (s *AccountsService) Delete(td *TokenData, id string) (err error) {

	// Get existing member
	var existing Entity
	if existing, err = s.sh.Database.Get(NewAccount, id); err != nil {
		return s.serviceError("Delete", err)
	}

	if existing.(*Account).Flag < 0 {
		if err = s.sh.Database.Delete(NewAccount, id); err != nil {
			return s.serviceError("Delete", err)
		} else {
			s.auditLog(td, existing, actionDelete, existing, nil)
			return nil
		}
	} else {
		existing.(*Account).Flag = -1
		existing.(*Account).Status = AccountStatusCodes.SUSPENDED

		if _, err = s.sh.Database.Update(existing); err != nil {
			return s.serviceError("Delete", err)
		} else {
			s.auditLog(td, existing, actionDelete, existing, nil)
			return nil
		}
	}
}

// Get single account by id
func (s *AccountsService) Get(td *TokenData, id string) (Entity, error) {

	if ent, err := s.sh.Database.Get(NewAccount, id); err != nil {
		return nil, fmt.Errorf("[%s]::Get: %v", s.ServiceName, err)
	} else {
		return ent, nil
	}
}

// AccountsFindParams Query params aggregator for find commands service
type AccountsFindParams struct {
	Search string              // Filter by free text search (using * wildcard)
	Status []AccountStatusCode // by status(s)
	Sort   string              // Sort descriptor (field name with suffix +/- for sort order)
	Page   int                 // Page number for pagination
	Size   int                 // Page size: number of items per page
}

func (f *AccountsFindParams) Statuses() (result []any) {
	for _, t := range f.Status {
		result = append(result, t)
	}
	return result
}

// Find list of accounts by filter
func (s *AccountsService) Find(p AccountsFindParams) (entities []Entity, total int64, pages int, error error) {
	if entities, total, error = s.sh.Database.Query(NewAccount).
		MatchAny(
			F("name").Like(p.Search),
			F("enName").Like(p.Search),
			F("email").Like(p.Search),
		).
		MatchAll(
			F("flag").Gte(0),
			F("status").In(p.Statuses()...),
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
