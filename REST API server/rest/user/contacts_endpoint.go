package rest

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	. "github.com/go-yaaf/yaaf-common/entity"
	"github.com/go-yaaf/yaaf-common/rest"

	. "github.com/go-yaaf/yaaf-examples/rest-api/model/entities"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/enums"
	. "github.com/go-yaaf/yaaf-examples/rest-api/rest"
	s "github.com/go-yaaf/yaaf-examples/rest-api/services"
)

// region Endpoint structure and factory method ------------------------------------------------------------------------

// ContactsEndPoint Services for contacts actions
// @Service: ContactsService
// @Path: /contacts
// @Context: usr-contacts
// @RequestHeader: X-API-KEY     | The key to identify the application (dashboard)
// @RequestHeader: Authorization | The bearer token to identify the logged-in user
// @ResourceGroup: Contacts Actions
type ContactsEndPoint struct {
	BaseEndPoint
	service *s.ContactsService
}

// NewContactsEndPoint factory method
func NewContactsEndPoint(service *s.ContactsService) RestEndpoint {
	return &ContactsEndPoint{service: service}
}

func (h *ContactsEndPoint) Path() string {
	return usrApiVersion + "/contacts"
}

func (h *ContactsEndPoint) RestEntries() (restEntries []RestEntry) {
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

// Get new and empty contact template
// @Http: POST /new
// @Return: EntityResponse<Contact>
func (h *ContactsEndPoint) new(c *gin.Context) {

	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	// Create empty entity
	entity := NewContact()
	entity.(*Contact).Id = ""
	entity.(*Contact).CreatedOn = 0
	entity.(*Contact).UpdatedOn = 0
	entity.(*Contact).Props = make(Json)
	c.JSON(http.StatusOK, rest.NewEntityResponse(entity))
}

// Create new contact
// @Http: POST /
// @BodyParam: body | Contact | contact data to create
// @Return: EntityResponse<Contact>
func (h *ContactsEndPoint) create(c *gin.Context) {

	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	// Read entity from body
	entity := NewContact()
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

// Update existing contact
// @Http: PUT /
// @BodyParam: body | Contact | contact data to update
// @Return: EntityResponse<Contact>
func (h *ContactsEndPoint) update(c *gin.Context) {
	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	// Read entity from body
	entity := NewContact()
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

// Delete contact and all its content
// @Http: DELETE /{id}
// @PathParam: id | string | contact ID to delete
// @Return: ActionResponse
func (h *ContactsEndPoint) delete(c *gin.Context) {
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

// Get a single contact by id
// @Http: GET /{id}
// @PathParam: id | string | contact ID to fetch
// @Return: EntityResponse<Contact>
func (h *ContactsEndPoint) get(c *gin.Context) {
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

// Find contacts by query
// @Http: GET /
// @QueryParam: search | string              | filter contacts by free text search
// @QueryParam: status | []StatusCode        | filter contacts by status(es)
// @QueryParam: sort   | string              | sort results by field and direction: (e.g. time = sort by time asc, time- = sort by time desc)
// @QueryParam: page   | int                 | page number (for pagination)
// @QueryParam: size   | int                 | number of items per page (for pagination)
// @Return: EntitiesResponse<Contact>
func (h *ContactsEndPoint) find(c *gin.Context) {
	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	p := s.ContactsFindParams{
		Search: h.GetParamAsString(c, "search", ""),
		Status: h.GetParamAsEnumArray(c, "status", *StatusCodes),
		Sort:   h.GetParamAsString(c, "sort", "lastName"),
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
