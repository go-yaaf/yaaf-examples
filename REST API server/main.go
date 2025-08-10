package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-yaaf/yaaf-common/logger"
	"github.com/go-yaaf/yaaf-examples/rest-api/cmd/server"
	"github.com/go-yaaf/yaaf-examples/rest-api/common"
	"github.com/go-yaaf/yaaf-examples/rest-api/config"
	"github.com/go-yaaf/yaaf-examples/rest-api/rest"
	usr "github.com/go-yaaf/yaaf-examples/rest-api/rest/user"
)

func init() {
	logLevel := config.GetConfig().LogLevel()
	logger.SetLevel(logLevel)
	logger.EnableJsonFormat(config.GetConfig().EnableLogJsonFormat())
	logger.Init()
}

// Main entry point
func main() {

	if app, err := setup(); err != nil {
		logger.Error("bootstrap error: %s", err.Error())
	} else {
		logger.Info("Starting the service...")
		app.Start()
	}
}

// Bootstrap the application components and inject all dependencies
func setup() (*server.Application, error) {

	// Init config
	serviceConfig := config.GetConfig()

	// Init service hub
	facade := common.NewServiceHub()

	// Init REST server
	restServer := newRestServer(serviceConfig, facade)

	// Init application
	application, err := server.NewApplication(serviceConfig, facade, restServer)
	if err != nil {
		return nil, err
	}
	return application, nil
}

// Create the REST server for the prometheus metrics endpoint
func newRestServer(cfg *config.ServiceConfig, facade *common.ServiceHub) *rest.Server {
	restServer := rest.NewRESTServer(cfg)
	restServer.AddEndpoints(usr.NewListOfUserRestEndPoints(facade)...)
	restServer.AddEndpoints(rest.NewHealthEndPoint())

	// Add documentation endpoint
	restServer.AddStaticEndpoint("/doc", "./doc")

	return restServer
}

// Log directory structure, for DEBUG only
// Motivation: to understand volume mapping of Google bucket to Google Cloud Run container
func logFileSystemStructure(root string) {
	fmt.Println("Reading content of:", root)
	files, err := os.ReadDir(root)
	if err != nil {
		fmt.Println("Error reading dir:", err)
		return
	}
	fmt.Println("Total files in ", root, " :", len(files))
	for i, f := range files {
		if i > 10 {
			break
		}
		fmt.Println("Found file:", f.Name())

		if absPath, er := filepath.Abs(filepath.Join(root, f.Name())); er != nil {
			fmt.Println("Error opening file path:", er)
		} else {
			fmt.Println(absPath)
		}
	}
	fmt.Println("--------------------")

}
