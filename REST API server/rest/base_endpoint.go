package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	mc "github.com/go-yaaf/yaaf-examples/rest-api/model/common"
	"github.com/go-yaaf/yaaf-examples/rest-api/utils"
)

// RestEntry represent a single HTTP REST call
type RestEntry struct {
	Path, // Rest method path
	Method string // HTTP method verb
	Handler gin.HandlerFunc // Handler function
}

// RestEndpoint is a group of RestEntry
type RestEndpoint interface {
	Path() string             // Rest method path
	RestEntries() []RestEntry // List of REST entries
}

type BaseEndPoint struct{}

// GetTokenData extract security token data from Authorization header
func (b *BaseEndPoint) GetTokenData(c *gin.Context) *mc.TokenData {

	token := c.GetHeader("X-ACCESS-TOKEN")
	if td, err := utils.TokenUtils().ParseToken(token); err != nil {
		_ = c.AbortWithError(http.StatusForbidden, fmt.Errorf("invalid access token"))
		return nil
	} else {
		return td
	}
}

// GetTimezoneOffset returns the value of timezone offset header in minutes
func (b *BaseEndPoint) GetTimezoneOffset(c *gin.Context) int {

	now := time.Now()
	_, offsetSeconds := now.Zone()
	localOffset := offsetSeconds / 60

	offset := c.GetHeader("X-TIMEZONE-OFFSET")
	clientOffset, err := strconv.Atoi(offset)
	if err != nil {
		return 0
	} else {
		return localOffset + clientOffset
	}
}
