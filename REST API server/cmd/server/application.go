// Package server The application hosts all the facilities (services)
package server

import (
	"fmt"
	"time"

	"github.com/go-yaaf/yaaf-common/logger"
	"github.com/go-yaaf/yaaf-examples/rest-api/common"
	"github.com/go-yaaf/yaaf-examples/rest-api/config"
	"github.com/go-yaaf/yaaf-examples/rest-api/rest"
)

// region Application structure and factory method ---------------------------------------------------------------------

// Application is the main server struct
// It contains and manages the state of the application and apply business login
type Application struct {
	config *config.ServiceConfig // Application configuration
	server *rest.Server          // REST server
	facade *common.ServiceHub    // Group all facility services (database, elastic, streaming etc)
}

// NewApplication Factory method
func NewApplication(cfg *config.ServiceConfig, hub *common.ServiceHub, server *rest.Server) (*Application, error) {
	return &Application{
		config: cfg,
		server: server,
		facade: hub,
	}, nil
}

// endregion

// region Application methods ------------------------------------------------------------------------------------------

// Start the application
func (app *Application) Start() {

	//application = app

	// Verify database schema
	if err := verifyDatabaseSchema(app.facade.Database); err != nil {
		logger.Fatal("error initializing database: %s", err.Error())
		return
	}

	// Verify root account schema
	if err := initializeDatabase(app.facade.Database); err != nil {
		logger.Fatal("error initializing database: %s", err.Error())
		return
	}

	// Start REST server for prometheus metrics endpoint
	go func() {
		app.startRestServer()
		logger.Info("Closing the REST server...")
	}()

	<-make(chan struct{})
}

// Start the REST server
func (app *Application) startRestServer() {

	port := config.GetConfig().ServerPort()
	if port == 0 {
		return
	}

	logger.Info(logTimezoneOffset())
	logger.Info("Starting REST server, listening on port: %d", port)

	if err := app.server.Start(port); err != nil {
		logger.Warn("error starting REST server: %s", err.Error())
	}
}

func logTimezoneOffset() string {
	now := time.Now()
	zone, offsetSeconds := now.Zone()
	localOffset := offsetSeconds / 60
	return fmt.Sprintf("Local timezone: %s, UTC offset: %d", zone, localOffset)
}

// endregion
