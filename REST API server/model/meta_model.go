package model

import (
	"strings"

	. "github.com/go-yaaf/yaaf-common/entity"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/entities"
)

var EntityRepo map[string]EntityFactory = make(map[string]EntityFactory)

// RegisterEntity register entity factory
func registerEntity(ef EntityFactory) {
	EntityRepo[ef().TABLE()] = ef
}

// GetEntityFactory returns entity factory by name
func GetEntityFactory(name string) EntityFactory {
	name = strings.ToLower(name)
	if ef, exists := EntityRepo[name]; exists {
		return ef
	} else {
		return nil
	}
}

func init() {
	registerEntity(NewAccount)
	registerEntity(NewAuditLog)
	registerEntity(NewContact)
	registerEntity(NewUser)
	registerEntity(NewUsersGroup)
}
