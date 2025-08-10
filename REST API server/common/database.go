package common

import (
	"strings"

	pgsql "github.com/go-yaaf/yaaf-common-postgresql/postgresql"
	"github.com/go-yaaf/yaaf-common/database"

	"github.com/go-yaaf/yaaf-examples/rest-api/config"
)

// NewDatabase is the factory method for a concrete implementation of the IDatabase interface
// In this project we support two implementations: in-memory database (for testing) and postgresql (for production)
// The concrete implementation is defined by the database URI schema
func NewDatabase() database.IDatabase {

	uri := config.GetConfig().DatabaseUri()

	// For postgresql schema, create postgres
	if strings.HasPrefix(uri, "postgres://") {
		if db, err := pgsql.NewPostgresDatabase(uri); err != nil {
			panic(err)
		} else {
			return db
		}
	}

	// For unknown or empty schema, create local in-memory DB
	db, err := database.NewInMemoryDatabase()
	if err != nil {
		panic(err)
	} else {
		return db
	}
}
