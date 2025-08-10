package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-yaaf/yaaf-common/rest"
	mc "github.com/go-yaaf/yaaf-examples/rest-api/model/common"
	. "github.com/go-yaaf/yaaf-examples/rest-api/rest"
	s "github.com/go-yaaf/yaaf-examples/rest-api/services"
	"net/http"
	"sort"
)

// region Endpoint structure and factory method ------------------------------------------------------------------------

// UserEndPoint Services for user registration and login
// @Service: UserService
// @Path: /user
// @Context: usr-user
// @RequestHeader: X-API-KEY     | The key to identify the application (dashboard)
// @RequestHeader: Authorization | The bearer token to identify the logged-in user
// @ResourceGroup: User Actions
type UserEndPoint struct {
	BaseEndPoint
	service *s.UsersService
}

// NewUserEndPoint factory method
func NewUserEndPoint(service *s.UsersService) RestEndpoint {
	return &UserEndPoint{service: service}
}

func (h *UserEndPoint) Path() string {
	return usrApiVersion + "/user"
}

func (h *UserEndPoint) RestEntries() (restEntries []RestEntry) {
	restEntries = []RestEntry{
		{Method: http.MethodPost, Handler: h.authorize, Path: "/authorize"},
		// {Method: http.MethodGet, Handler: h.enums, Path: "/enums"},
	}

	// Sort entries for best match
	sort.Slice(restEntries, func(i, j int) bool {
		return restEntries[i].Path > restEntries[j].Path
	})
	return
}

// endregion

// region Endpoint REST handlers ---------------------------------------------------------------------------------------

// Authorize user, verify user exists in the system
// The response includes access token valid for 20 minutes. The client side should renew the token before expiration using refresh-token method
// @Http: POST /authorize
// @BodyParam: body | LoginParams | User verified email
// @Return: EntityResponse<User>
func (h *UserEndPoint) authorize(c *gin.Context) {

	// Read email from body
	login := mc.LoginParams{}
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusUnauthorized, rest.NewErrorResponse(errors.New("unauthorized")))
		return
	}

	if user, token, err := h.service.Authorize(login.Email); err != nil {
		c.JSON(http.StatusUnauthorized, rest.NewErrorResponse(errors.New("unauthorized")))
		return
	} else {
		c.Header("X-ACCESS-TOKEN", token)
		c.JSON(http.StatusOK, rest.NewEntityResponse(user))
	}
}

// endregion
