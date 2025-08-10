package common

import (
	"strings"

	rds "github.com/go-yaaf/yaaf-common-redis/redis"
	"github.com/go-yaaf/yaaf-common/database"

	"github.com/go-yaaf/yaaf-examples/rest-api/config"
)

// NewDataCache is the factory method for a concrete implementation of the IDataCache interface
// In this project we support two implementations: in-memory cache (for testing) and Redis (for production)
// The concrete implementation is defined by the datacache URI schema
func NewDataCache() database.IDataCache {

	uri := config.GetConfig().DataCacheUri()

	// For redis schema, create redis implementation
	if strings.HasPrefix(uri, "redis://") {
		if dc, err := rds.NewRedisDataCache(uri); err != nil {
			panic(err)
		} else {
			return dc
		}
		return nil
	}

	// For unknown or empty schema, create local in-memory DB
	dc, err := database.NewInMemoryDataCache()
	if err != nil {
		panic(err)
	} else {
		return dc
	}
}
