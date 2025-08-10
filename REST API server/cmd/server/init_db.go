package server

import (
	"fmt"
	"strings"

	"github.com/go-yaaf/yaaf-common/database"
	"github.com/go-yaaf/yaaf-common/logger"

	. "github.com/go-yaaf/yaaf-examples/rest-api/config"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/entities"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/enums"
)

// verifyDatabaseSchema create all the tables and indexes
func verifyDatabaseSchema(database database.IDatabase) error {

	// Fill table and indexes
	ddl := make(map[string][]string)

	ddl["account"] = []string{"name", "status", "flag"}
	ddl["audit_log"] = []string{"createdOn", "accountId", "userId", "action", "itemType", "itemId", "itemName"}
	ddl["contact"] = []string{"firstName", "lastName", "status", "updatedOn", "flag"}
	ddl["user"] = []string{"name", "email"}
	ddl["users_group"] = []string{"name", "updatedOn"}

	if err := database.ExecuteDDL(ddl); err != nil {
		return err
	}

	return nil
}

// initializeDatabase create the initial configuration data
func initializeDatabase(database database.IDatabase) error {
	if err := verifyRootAdmin(database); err != nil {
		return err
	}
	return nil
}

// verifyRootAdmin create the initial system administrator if not exists
func verifyRootAdmin(database database.IDatabase) error {
	email := GetConfig().InitialAdminEmail()
	if len(email) == 0 {
		return nil
	}

	user := NewUser()

	if idx := strings.Index(email, "@"); idx > -1 {
		user.(*User).Id = email
		user.(*User).Email = email
		user.(*User).Name = email[:idx]

	} else {
		return fmt.Errorf("%s not a valid email", email)
	}

	user.(*User).Type = UserTypeCodes.SYSADMIN
	user.(*User).Status = UserStatusCodes.ACTIVE

	if exists, err := database.Exists(NewUser, user.ID()); err != nil {
		return err
	} else {
		if exists {
			return nil
		}
		if added, er := database.Insert(user); er != nil {
			return er
		} else {
			logger.Info("root admin created: %s", added.ID())
			return nil
		}
	}
}
