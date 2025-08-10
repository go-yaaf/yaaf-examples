package services

import (
	"errors"
	"fmt"
	"sync"

	. "github.com/go-yaaf/yaaf-common/database"
	. "github.com/go-yaaf/yaaf-common/entity"

	"github.com/go-yaaf/yaaf-examples/rest-api/common"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/common"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/entities"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/enums"
	. "github.com/go-yaaf/yaaf-examples/rest-api/utils"
)

var usersServiceOnce sync.Once
var usersServiceInst *UsersService = nil

type UsersService struct {
	BaseService
	sh *common.ServiceHub // Service hub
}

// GetUsersService factory function
func GetUsersService(sh *common.ServiceHub) *UsersService {
	usersServiceOnce.Do(func() {
		if usersServiceInst == nil {
			usersServiceInst = &UsersService{
				BaseService: BaseService{ServiceName: "UsersService"},
				sh:          sh,
			}
		}
	})
	return usersServiceInst
}

// Create new user in the system
func (s *UsersService) Create(td *TokenData, entity Entity) (Entity, error) {

	ent := entity.(*User)

	if StringUtils().IsValidEmail(ent.Email) {
		ent.Id = ent.Email
	} else {
		return nil, s.serviceError("Create", errors.New("not a valid email"))
	}

	// Override system fields,
	if len(ent.Id) == 0 {
		ent.Id = TokenUtils().ShortID()
	}
	ent.CreatedOn = Now()
	ent.UpdatedOn = Now()
	ent.Props = nil

	if updated, er := s.sh.Database.Insert(ent); er != nil {
		return nil, s.serviceError("Create", er)
	} else {
		s.auditLog(td, ent, actionCreate, nil, updated)
		return updated, nil
	}
}

// Update existing user in the system
func (s *UsersService) Update(td *TokenData, entity Entity) (Entity, error) {

	ent := entity.(*User)

	// Get existing account
	existing, err := s.sh.Database.Get(NewUser, ent.Id)
	if err != nil {
		return nil, s.serviceError("Update", err)
	}

	// Override system fields,
	ent.CreatedOn = existing.(*User).CreatedOn
	ent.UpdatedOn = Now()
	ent.Props = nil

	if updated, er := s.sh.Database.Update(ent); er != nil {
		return nil, s.serviceError("Update", er)
	} else {
		s.auditLog(td, ent, actionUpdate, existing, updated)
		return updated, nil
	}
}

// Delete user
func (s *UsersService) Delete(td *TokenData, id string) (err error) {

	// Get existing member
	var existing Entity
	if existing, err = s.sh.Database.Get(NewUser, id); err != nil {
		return s.serviceError("Delete", err)
	}

	if existing.(*User).Flag < 0 {
		if err = s.sh.Database.Delete(NewUser, id); err != nil {
			return s.serviceError("Delete", err)
		} else {
			s.auditLog(td, existing, actionDelete, existing, nil)
			return nil
		}
	} else {
		existing.(*User).Flag = -1

		if _, err = s.sh.Database.Update(existing); err != nil {
			return s.serviceError("Delete", err)
		} else {
			s.auditLog(td, existing, actionDelete, existing, nil)
			return nil
		}
	}
}

// Get a single user by id
func (s *UsersService) Get(td *TokenData, id string) (Entity, error) {
	if ent, err := s.sh.Database.Get(NewUser, id); err != nil {
		return nil, fmt.Errorf("[%s]::Get: %v", s.ServiceName, err)
	} else {
		return ent, nil
	}
}

// GetBtEmail get single user by email
func (s *UsersService) GetBtEmail(email string) (Entity, error) {
	if ent, err := s.sh.Database.Query(NewUser).
		Filter(F("email").Eq(email)).
		FindSingle(); err != nil {
		return nil, s.serviceError("GetBtEmail", err)
	} else {
		return ent, err
	}
}

// Authorize get a single user by email, get the member of account and create JWT token
func (s *UsersService) Authorize(email string) (user Entity, token string, error error) {
	// Get user by email
	user, error = s.sh.Database.Query(NewUser).Filter(F("email").Eq(email)).FindSingle()
	if error != nil {
		return nil, "", s.serviceError("Authorize", error)
	}

	if user.(*User).Status != UserStatusCodes.ACTIVE {
		return nil, "", s.serviceError("Authorize", fmt.Errorf("not authorized"))
	}

	// Update last sign-in
	user.(*User).LastSignIn = Now()
	_, _ = s.sh.Database.Update(user)

	// if user is a sysadmin, return default account
	if user.(*User).Type == UserTypeCodes.SYSADMIN {
		token, error = s.createTokenForSysAdmin(user)
		user.(*User).Roles = UserRoleFlags.ALL
	} else {
		token, error = s.createTokenForUser(user)
	}
	return
}

// UsersFindParams Query params aggregator for find commands service
type UsersFindParams struct {
	Search string           // Filter by text search on name / id / email (using * wildcard)
	Type   []UserTypeCode   // Filter by type(s)
	Status []UserStatusCode // Filter by status(s)
	Sort   string           // Sort descriptor (field name with suffix +/- for sort order)
	Page   int              // Page number for pagination
	Size   int              // Page size: number of items per page
}

// Find a list of members by filter
func (s *UsersService) Find(p UsersFindParams) (entities []Entity, total int64, pages int, error error) {
	if entities, total, error = s.sh.Database.Query(NewUser).
		MatchAny(
			F("id").Eq(p.Search),
			F("name").Like(p.Search),
			F("email").Like(p.Search),
		).
		MatchAll(
			F("flag").Gte(0),
			F("type").In(ToAnyVariadic(p.Type)...),
			F("status").In(ToAnyVariadic(p.Status)...),
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

// Create token for sys admin
func (s *UsersService) createTokenForSysAdmin(user Entity) (string, error) {
	td := &TokenData{
		SubjectId:   user.ID(),
		SubjectType: UserTypeCodes.SYSADMIN,
		Status:      user.(*User).Status,
		ExpiresIn:   int64(Now() + 1000*60*30),
	}
	// Update default account
	if token, err := TokenUtils().CreateToken(td); err != nil {
		return "", s.serviceError("createTokenForSysAdmin", err)
	} else {
		return token, nil
	}
}

// Create token for user
func (s *UsersService) createTokenForUser(user Entity) (string, error) {
	td := &TokenData{
		SubjectId:   user.ID(),
		SubjectType: UserTypeCodes.SYSADMIN,
		Status:      user.(*User).Status,
		ExpiresIn:   int64(Now() + 1000*60*30),
	}

	// Update default account
	if token, err := TokenUtils().CreateToken(td); err != nil {
		return "", s.serviceError("createTokenForUser", err)
	} else {
		return token, nil
	}
}
