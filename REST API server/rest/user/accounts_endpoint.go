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

// AccountsEndPoint Services for user registration and login
// @Service: AccountsService
// @Path: /accounts
// @Context: usr-accounts
// @RequestHeader: X-API-KEY     | The key to identify the application (dashboard)
// @RequestHeader: Authorization | The bearer token to identify the logged-in user
// @ResourceGroup: Accounts Actions
type AccountsEndPoint struct {
	BaseEndPoint
	service *s.AccountsService
}

// NewAccountsEndPoint factory method
func NewAccountsEndPoint(service *s.AccountsService) RestEndpoint {
	return &AccountsEndPoint{service: service}
}

func (h *AccountsEndPoint) Path() string {
	return usrApiVersion + "/accounts"
}

func (h *AccountsEndPoint) RestEntries() (restEntries []RestEntry) {
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

// Get new and empty account template
// @Http: POST /new
// @Return: EntityResponse<Account>
func (h *AccountsEndPoint) new(c *gin.Context) {

	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	// Create empty entity
	entity := NewAccount()
	entity.(*Account).Id = ""
	entity.(*Account).CreatedOn = 0
	entity.(*Account).UpdatedOn = 0
	entity.(*Account).Props = make(Json)
	entity.(*Account).Status = AccountStatusCodes.ACTIVE

	c.JSON(http.StatusOK, rest.NewEntityResponse(entity))
}

// Create new account
// @Http: POST /
// @BodyParam: body | Account | account data to create
// @Return: EntityResponse<Account>
func (h *AccountsEndPoint) create(c *gin.Context) {

	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	// Read entity from body
	entity := NewAccount()
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

// Update existing account
// @Http: PUT /
// @BodyParam: body | Account | account data to update
// @Return: EntityResponse<Account>
func (h *AccountsEndPoint) update(c *gin.Context) {
	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	// Read entity from body
	entity := NewAccount()
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

// Delete account and all its content
// @Http: DELETE /{id}
// @PathParam: id | string | account ID to delete
// @Return: ActionResponse
func (h *AccountsEndPoint) delete(c *gin.Context) {
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

// Get a single account by id
// @Http: GET /{id}
// @PathParam: id | string | account ID to fetch
// @Return: EntityResponse<Account>
func (h *AccountsEndPoint) get(c *gin.Context) {
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

// Find accounts by query
// @Http: GET /
// @QueryParam: search | string              | filter accounts by free text search on account id, account name
// @QueryParam: status | []AccountStatusCode | filter accounts by status(s)
// @QueryParam: sort   | string              | sort results by field and direction: (e.g. time = sort by time asc, time- = sort by time desc)
// @QueryParam: page   | int                 | page number (for pagination)
// @QueryParam: size   | int                 | number of items per page (for pagination)
// @Return: EntitiesResponse<Account>
func (h *AccountsEndPoint) find(c *gin.Context) {
	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	p := s.AccountsFindParams{
		Search: h.GetParamAsString(c, "search", ""),
		Status: h.GetParamAsEnumArray(c, "status", *AccountStatusCodes),
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
