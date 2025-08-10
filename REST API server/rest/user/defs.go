package rest

import (
	"github.com/go-yaaf/yaaf-examples/rest-api/common"
	"github.com/go-yaaf/yaaf-examples/rest-api/rest"
	s "github.com/go-yaaf/yaaf-examples/rest-api/services"
)

const (
	usrApiVersion = "/v1"
)

// NewListOfUserRestEndPoints is a factory method for user endpoints list
func NewListOfUserRestEndPoints(facade *common.ServiceHub) []rest.RestEndpoint {
	list := make([]rest.RestEndpoint, 0)
	list = append(list, NewAccountsEndPoint(s.GetAccountsService(facade)))
	list = append(list, NewAuditLogsEndPoint(s.GetAuditLogsService(facade)))
	list = append(list, NewContactsEndPoint(s.GetContactsService(facade)))
	list = append(list, NewGroupsEndPoint(s.GetGroupsService(facade)))
	list = append(list, NewUserEndPoint(s.GetUsersService(facade)))
	list = append(list, NewUsersEndPoint(s.GetUsersService(facade)))

	return list
}

//func createLookup(list []Entity) []Entity {
//	lookup := make([]Entity, 0)
//
//	for _, ent := range list {
//		se := &Tuple[string, string]{
//			Key:   ent.ID(),
//			Value: ent.NAME(),
//		}
//		lookup = append(lookup, se)
//	}
//
//	return lookup
//}
