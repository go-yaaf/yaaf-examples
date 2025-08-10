package common

import (
	"os"
	"strings"

	"github.com/go-yaaf/yaaf-common/database"
)

// ServiceHub is the main application hub for all middleware facilities (e.g. database, cache, messaging, etc`)
type ServiceHub struct {
	Database  database.IDatabase  // Configuration database middleware facade
	DataCache database.IDataCache // Distributed cache middleware facade
	Version   string              // Current service version
}

var facade *ServiceHub = nil

// SetServiceHub set the service hub instance
func SetServiceHub(sh *ServiceHub) {
	facade = sh
}

// GetServiceHub get the service hub instance
func GetServiceHub() *ServiceHub {
	return facade
}

// NewServiceHub is a service hub factory method
func NewServiceHub() *ServiceHub {
	facade = &ServiceHub{
		Database:  NewDatabase(),
		DataCache: NewDataCache(),
		Version:   getVersion(),
	}

	return facade
}

func getVersion() string {
	version := "1.0.0"

	// Try to read build tag from current folder
	if versionBytes, err := os.ReadFile("build-tag"); err == nil {
		version = string(versionBytes)
		if strings.Contains(version, "\n") {
			n := strings.Index(version, "\n")
			version = version[:n]
		}
	}
	return version
}
