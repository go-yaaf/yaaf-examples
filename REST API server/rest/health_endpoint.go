package rest

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-yaaf/yaaf-common/rest"
)

// region Endpoint structure and factory method ------------------------------------------------------------------------

// HealthEndPoint for health check
type HealthEndPoint struct {
	BaseEndPoint
}

// NewHealthEndPoint factory method
func NewHealthEndPoint() RestEndpoint {
	return &HealthEndPoint{}
}

// endregion

// region Endpoint methods implementation ------------------------------------------------------------------------------

func (h *HealthEndPoint) Path() string {
	return "/"
}

func (h *HealthEndPoint) RestEntries() (restEntries []RestEntry) {
	restEntries = []RestEntry{
		{Method: http.MethodGet, Handler: h.root, Path: "/"},
	}
	return
}

// Root handler returns the current version number of the API
func (h *HealthEndPoint) root(c *gin.Context) {

	version := "1.0.0"

	// Try to read build tag from current folder
	if versionBytes, err := os.ReadFile("build-tag"); err == nil {
		version = string(versionBytes)
		if strings.Contains(version, "\n") {
			n := strings.Index(version, "\n")
			version = version[:n]
		}
	}

	// Read build number
	c.JSON(http.StatusOK, rest.NewActionResponse("rest-api", version))
}

// endregion
