package rest

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/go-yaaf/yaaf-common/rest"

	. "github.com/go-yaaf/yaaf-examples/rest-api/model/entities"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/enums"
	. "github.com/go-yaaf/yaaf-examples/rest-api/rest"
	s "github.com/go-yaaf/yaaf-examples/rest-api/services"
)

// region Endpoint structure and factory method ------------------------------------------------------------------------

// UsersEndPoint Services for users actions
// @Service: UsersService
// @Path: /users
// @Context: usr-users
// @RequestHeader: X-API-KEY     | The key to identify the application (dashboard)
// @RequestHeader: Authorization | The bearer token to identify the logged-in user
// @ResourceGroup: Users Actions
type UsersEndPoint struct {
	BaseEndPoint
	service *s.UsersService
}

// NewUsersEndPoint factory method
func NewUsersEndPoint(service *s.UsersService) RestEndpoint {
	return &UsersEndPoint{service: service}
}

func (h *UsersEndPoint) Path() string {
	return usrApiVersion + "/users"
}

func (h *UsersEndPoint) RestEntries() (restEntries []RestEntry) {
	restEntries = []RestEntry{
		{Method: http.MethodPost, Handler: h.create, Path: ""},
		{Method: http.MethodPost, Handler: h.create, Path: "/"},
		{Method: http.MethodPost, Handler: h.new, Path: "/new"},

		{Method: http.MethodPut, Handler: h.update, Path: ""},
		{Method: http.MethodPut, Handler: h.update, Path: "/"},

		{Method: http.MethodDelete, Handler: h.delete, Path: "/:id"},
		{Method: http.MethodGet, Handler: h.get, Path: "/:id"},

		{Method: http.MethodGet, Handler: h.find, Path: ""},
		{Method: http.MethodGet, Handler: h.find, Path: "/"},
	}

	// Sort entries for best match
	sort.Slice(restEntries, func(i, j int) bool {
		return restEntries[i].Path > restEntries[j].Path
	})
	return
}

// endregion

// region Endpoint REST handlers ---------------------------------------------------------------------------------------

// Get new and empty user template
// @Http: POST /new
// @Return: EntityResponse<User>
func (h *UsersEndPoint) new(c *gin.Context) {

	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	// Create empty entity
	entity := NewUser()
	entity.(*User).Id = ""
	entity.(*User).CreatedOn = 0
	entity.(*User).UpdatedOn = 0
	entity.(*User).Status = UserStatusCodes.PENDING
	entity.(*User).Type = UserTypeCodes.USER
	entity.(*User).Roles = UserRoleFlags.MANAGER

	c.JSON(http.StatusOK, rest.NewEntityResponse(entity))
}

// Create new user
// @Http: POST /
// @BodyParam: body | User | user data to create
// @Return: EntityResponse<User>
func (h *UsersEndPoint) create(c *gin.Context) {

	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	// Read entity from body
	entity := NewUser()
	if err := c.ShouldBindJSON(entity); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if result, err := h.service.Create(td, entity); err != nil {
		c.JSON(http.StatusInternalServerError, rest.NewErrorResponse(err))
	} else {
		c.JSON(http.StatusOK, rest.NewEntityResponse(result))
	}
}

// Update existing user
// @Http: PUT /
// @BodyParam: body | User | user data to update
// @Return: EntityResponse<User>
func (h *UsersEndPoint) update(c *gin.Context) {
	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	// Read entity from body
	entity := NewUser()
	if err := c.ShouldBindJSON(entity); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if result, err := h.service.Update(td, entity); err != nil {
		c.JSON(http.StatusInternalServerError, rest.NewErrorResponse(err))
	} else {
		c.JSON(http.StatusOK, rest.NewEntityResponse(result))
	}
}

// Delete user and all its content
// @Http: DELETE /{id}
// @PathParam: id | string | user ID to delete
// @Return: ActionResponse
func (h *UsersEndPoint) delete(c *gin.Context) {
	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	// Only Admin can delete this
	if td.SubjectType != UserTypeCodes.SYSADMIN {
		c.JSON(http.StatusInternalServerError, rest.NewErrorResponse(fmt.Errorf("delete is forbidden")))
		return
	}

	id := c.Params.ByName("id")

	if err := h.service.Delete(td, id); err != nil {
		c.JSON(http.StatusInternalServerError, rest.NewErrorResponse(err))
	} else {
		c.JSON(http.StatusOK, rest.NewActionResponse(td.SubjectId, id))
	}
}

// Get a single user by id
// @Http: GET /{id}
// @PathParam: id | string | user ID to fetch
// @Return: EntityResponse<User>
func (h *UsersEndPoint) get(c *gin.Context) {
	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	id := c.Params.ByName("id")

	if entity, err := h.service.Get(td, id); err != nil {
		c.JSON(http.StatusInternalServerError, rest.NewErrorResponse(err))
	} else {
		c.JSON(http.StatusOK, rest.NewEntityResponse(entity))
	}
}

// Find users by query
// @Http: GET /
// @QueryParam: search | string              | filter users by free text search
// @QueryParam: type   | []UserTypeCode      | filter users by type(s)
// @QueryParam: status | []UserStatusCode    | filter users by status(es)
// @QueryParam: sort   | string              | sort results by field and direction: (e.g. time = sort by time asc, time- = sort by time desc)
// @QueryParam: page   | int                 | page number (for pagination)
// @QueryParam: size   | int                 | number of items per page (for pagination)
// @Return: EntitiesResponse<User>
func (h *UsersEndPoint) find(c *gin.Context) {
	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	p := s.UsersFindParams{
		Search: h.GetParamAsString(c, "search", ""),
		Type:   h.GetParamAsEnumArray(c, "type", *UserTypeCodes),
		Status: h.GetParamAsEnumArray(c, "status", *UserStatusCodes),
		Sort:   h.GetParamAsString(c, "sort", "name"),
		Page:   h.GetParamAsInt(c, "page", 1),
		Size:   h.GetParamAsInt(c, "size", 100),
	}
	if list, total, _, err := h.service.Find(p); err != nil {
		c.JSON(http.StatusInternalServerError, rest.NewErrorResponse(err))
	} else {
		c.JSON(http.StatusOK, rest.NewEntitiesResponse(list, p.Page, p.Size, int(total)))
	}
}

// endregion
