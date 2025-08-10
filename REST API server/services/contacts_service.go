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

var contactsServiceOnce sync.Once
var contactsServiceInst *ContactsService = nil

type ContactsService struct {
	BaseService
	sh *ServiceHub // Service hub
}

// GetContactsService factory function
func GetContactsService(sh *ServiceHub) *ContactsService {
	contactsServiceOnce.Do(func() {
		if contactsServiceInst == nil {
			contactsServiceInst = &ContactsService{BaseService: BaseService{ServiceName: "ContactsService"}, sh: sh}
		}
	})
	return contactsServiceInst
}

// Create new contact in the system
func (s *ContactsService) Create(td *TokenData, entity Entity) (Entity, error) {

	ent := entity.(*Contact)

	// Override system fields,
	ent.Id = TokenUtils().GUID()
	ent.CreatedOn = Now()
	ent.UpdatedOn = Now()
	ent.Props = nil

	// Strip phone numbers
	ent.Mobile = s.stripPhone(ent.Mobile)

	if updated, er := s.sh.Database.Insert(ent); er != nil {
		return nil, s.serviceError("Create", er)
	} else {
		s.auditLog(td, ent, actionCreate, nil, updated)
		return updated, nil
	}
}

// Update existing contact in the system
func (s *ContactsService) Update(td *TokenData, entity Entity) (Entity, error) {

	ent := entity.(*Contact)

	// Get existing contact
	existing, err := s.sh.Database.Get(NewContact, ent.Id)
	if err != nil {
		return nil, s.serviceError("Update", err)
	}

	// Override system fields,
	ent.CreatedOn = existing.(*Contact).CreatedOn
	ent.UpdatedOn = Now()
	ent.Props = nil

	// Strip phone numbers
	ent.Mobile = s.stripPhone(ent.Mobile)

	if updated, er := s.sh.Database.Update(ent); er != nil {
		return nil, s.serviceError("Update", er)
	} else {
		s.auditLog(td, ent, actionUpdate, existing, updated)
		return updated, nil
	}
}

// Delete contact
func (s *ContactsService) Delete(td *TokenData, id string) (err error) {

	// Get existing contact
	var existing Entity
	if existing, err = s.sh.Database.Get(NewContact, id); err != nil {
		return s.serviceError("Delete", err)
	}

	if existing.(*Contact).Flag < 0 {
		if err = s.sh.Database.Delete(NewContact, id); err != nil {
			return s.serviceError("Delete", err)
		} else {
			s.auditLog(td, existing, actionDelete, existing, nil)
			return nil
		}
	} else {
		existing.(*Contact).Flag = -1

		if _, err = s.sh.Database.Update(existing); err != nil {
			return s.serviceError("Delete", err)
		} else {
			s.auditLog(td, existing, actionDelete, existing, nil)
			return nil
		}
	}
}

// Get single contact by id
func (s *ContactsService) Get(td *TokenData, id string) (Entity, error) {
	if ent, err := s.sh.Database.Get(NewContact, id); err != nil {
		return nil, fmt.Errorf("[%s]::Get: %v", s.ServiceName, err)
	} else {
		ent.(*Contact).Props = Json{}
		return ent, nil
	}
}

// ContactsFindParams Query params aggregator for find commands service
type ContactsFindParams struct {
	Search string       // Filter by free text search (using * wildcard)
	Status []StatusCode // Filter by status(es)
	Sort   string       // Sort descriptor (field name with suffix +/- for sort order)
	Page   int          // Page number for pagination
	Size   int          // Page size: number of items per page
}

// Find list of contacts by filter
func (s *ContactsService) Find(p ContactsFindParams) (entities []Entity, total int64, pages int, error error) {
	cb := func(in Entity) (out Entity) {
		in.(*Contact).Props = Json{}
		return in
	}

	if entities, total, error = s.sh.Database.Query(NewContact).
		MatchAny(
			F("id").Eq(p.Search),
			F("name").Like(p.Search),
			F("enName").Like(p.Search),
		).
		MatchAll(
			F("flag").Gte(0),
			F("status").In(ToAnyVariadic(p.Status)...),
		).
		Page(p.Page).
		Limit(p.Size).
		Sort(p.Sort).
		Apply(cb).
		Find(); error == nil {
		pages = s.calcPages(total, p.Size)
	} else {
		error = s.serviceError("Find", error)
	}
	return
}
