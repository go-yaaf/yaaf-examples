package rest

import (
	"fmt"
	"github.com/go-yaaf/yaaf-common/entity"
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

// GroupsEndPoint Services for groups actions
// @Service: GroupsService
// @Path: /groups
// @Context: usr-groups
// @RequestHeader: X-API-KEY     | The key to identify the application (dashboard)
// @RequestHeader: Authorization | The bearer token to identify the logged-in user
// @ResourceGroup: Groups Actions
type GroupsEndPoint struct {
	BaseEndPoint
	service *s.GroupsService
}

// NewGroupsEndPoint factory method
func NewGroupsEndPoint(service *s.GroupsService) RestEndpoint {
	return &GroupsEndPoint{service: service}
}

func (h *GroupsEndPoint) Path() string {
	return usrApiVersion + "/groups"
}

func (h *GroupsEndPoint) RestEntries() (restEntries []RestEntry) {
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

// Get new and empty flight template
// @Http: POST /new
// @Return: EntityResponse<UsersGroup>
func (h *GroupsEndPoint) new(c *gin.Context) {
	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	// Create empty entity
	result := NewUsersGroup()
	result.(*UsersGroup).Id = ""
	result.(*UsersGroup).CreatedOn = 0
	result.(*UsersGroup).UpdatedOn = 0
	result.(*UsersGroup).Props = make(entity.Json)
	result.(*UsersGroup).Name = ""
	result.(*UsersGroup).Email = ""

	c.JSON(http.StatusOK, rest.NewEntityResponse(result))
}

// Create new group
// @Http: POST /
// @BodyParam: body | Group | group data to create
// @Return: EntityResponse<Group>
func (h *GroupsEndPoint) create(c *gin.Context) {

	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	// Read entity from body
	ug := NewUsersGroup()
	if err := c.ShouldBindJSON(ug); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if result, err := h.service.Create(td, ug); err != nil {
		c.JSON(http.StatusInternalServerError, rest.NewErrorResponse(err))
	} else {
		c.JSON(http.StatusOK, rest.NewEntityResponse(result))
	}
}

// Update existing group
// @Http: PUT /
// @BodyParam: body | Group | group data to update
// @Return: EntityResponse<Group>
func (h *GroupsEndPoint) update(c *gin.Context) {
	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	// Read entity from body
	ug := NewUsersGroup()
	if err := c.ShouldBindJSON(ug); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if result, err := h.service.Update(td, ug); err != nil {
		c.JSON(http.StatusInternalServerError, rest.NewErrorResponse(err))
	} else {
		c.JSON(http.StatusOK, rest.NewEntityResponse(result))
	}
}

// Delete group and all its content
// @Http: DELETE /{id}
// @PathParam: id | string | group ID to delete
// @Return: ActionResponse
func (h *GroupsEndPoint) delete(c *gin.Context) {
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

// Get a single group by id
// @Http: GET /{id}
// @PathParam: id | string | group ID to fetch
// @Return: EntityResponse<Group>
func (h *GroupsEndPoint) get(c *gin.Context) {
	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	id := c.Params.ByName("id")

	if ent, err := h.service.Get(td, id); err != nil {
		c.JSON(http.StatusInternalServerError, rest.NewErrorResponse(err))
	} else {
		c.JSON(http.StatusOK, rest.NewEntityResponse(ent))
	}
}

// Find groups by query
// @Http: GET /
// @QueryParam: search | string              | filter groups by free text search
// @QueryParam: sort   | string              | sort results by field and direction: (e.g. time = sort by time asc, time- = sort by time desc)
// @QueryParam: page   | int                 | page number (for pagination)
// @QueryParam: size   | int                 | number of items per page (for pagination)
// @Return: EntitiesResponse<Group>
func (h *GroupsEndPoint) find(c *gin.Context) {
	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	p := s.GroupsFindParams{
		Search: h.GetParamAsString(c, "search", ""),
		Sort:   h.GetParamAsString(c, "sort", ""),
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
