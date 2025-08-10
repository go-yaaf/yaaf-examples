package rest

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/go-yaaf/yaaf-common/rest"

	. "github.com/go-yaaf/yaaf-examples/rest-api/model/entities"
	. "github.com/go-yaaf/yaaf-examples/rest-api/rest"
	s "github.com/go-yaaf/yaaf-examples/rest-api/services"
)

// region Endpoint structure and factory method ------------------------------------------------------------------------

// AuditLogsEndPoint Services for auditLogs actions
// @Service: AuditLogsService
// @Path: /audit_logs
// @Context: usr-auditLogs
// @RequestHeader: X-API-KEY     | The key to identify the application (dashboard)
// @RequestHeader: Authorization | The bearer token to identify the logged-in user
// @ResourceGroup: AuditLogs Actions
type AuditLogsEndPoint struct {
	BaseEndPoint
	service *s.AuditLogsService
}

// NewAuditLogsEndPoint factory method
func NewAuditLogsEndPoint(service *s.AuditLogsService) RestEndpoint {
	return &AuditLogsEndPoint{service: service}
}

func (h *AuditLogsEndPoint) Path() string {
	return usrApiVersion + "/audit_logs"
}

func (h *AuditLogsEndPoint) RestEntries() (restEntries []RestEntry) {
	restEntries = []RestEntry{
		{Method: http.MethodPost, Handler: h.create, Path: ""},
		{Method: http.MethodPost, Handler: h.create, Path: "/"},

		{Method: http.MethodGet, Handler: h.get, Path: "/:id"},

		{Method: http.MethodGet, Handler: h.find, Path: ""},
		{Method: http.MethodGet, Handler: h.find, Path: "/"},
		{Method: http.MethodGet, Handler: h.histogram, Path: "/histogram"},
	}

	// Sort entries for best match
	sort.Slice(restEntries, func(i, j int) bool {
		return restEntries[i].Path > restEntries[j].Path
	})
	return
}

// endregion

// region Endpoint REST handlers ---------------------------------------------------------------------------------------

// Create new auditLog
// @Http: POST /
// @BodyParam: body | AuditLog | auditLog data to create
// @Return: EntityResponse<AuditLog>
func (h *AuditLogsEndPoint) create(c *gin.Context) {

	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	// Read entity from body
	entity := NewAuditLog()
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

// Get a single auditLog by id
// @Http: GET /{id}
// @PathParam: id | string | auditLog ID to fetch
// @Return: EntityResponse<AuditLog>
func (h *AuditLogsEndPoint) get(c *gin.Context) {
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

// Find auditLogs by query
// @Http: GET /
// @QueryParam: from     | Timestamp           | start of time range filter
// @QueryParam: to       | Timestamp           | end of time range filter
// @QueryParam: userId   | string              | filter auditLogs by user id
// @QueryParam: action   | string              | filter auditLogs by action
// @QueryParam: itemType | string              | filter auditLogs by item type
// @QueryParam: itemId   | string              | filter auditLogs by item id
// @QueryParam: itemName | string              | filter auditLogs by item name
// @QueryParam: search   | string              | filter auditLogs by free text search
// @QueryParam: sort     | string              | sort results by field and direction: (e.g. time = sort by time asc, time- = sort by time desc)
// @QueryParam: page     | int                 | page number (for pagination)
// @QueryParam: size     | int                 | number of items per page (for pagination)
// @Return: EntitiesResponse<AuditLog>
func (h *AuditLogsEndPoint) find(c *gin.Context) {
	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	p := s.AuditLogsFindParams{
		From:     h.GetParamAsTimestamp(c, "from", 0),
		To:       h.GetParamAsTimestamp(c, "to", 0),
		UserId:   h.GetParamAsString(c, "userId", ""),
		Action:   h.GetParamAsString(c, "action", ""),
		ItemType: h.GetParamAsString(c, "itemType", ""),
		ItemId:   h.GetParamAsString(c, "itemId", ""),
		ItemName: h.GetParamAsString(c, "itemName", ""),
		Search:   h.GetParamAsString(c, "search", ""),
		Sort:     h.GetParamAsString(c, "sort", "createdOn-"),
		Page:     h.GetParamAsInt(c, "page", 1),
		Size:     h.GetParamAsInt(c, "size", 100),
	}

	if list, total, _, err := h.service.Find(p); err != nil {
		c.JSON(http.StatusInternalServerError, rest.NewErrorResponse(err))
	} else {
		c.JSON(http.StatusOK, rest.NewEntitiesResponse(list, p.Page, p.Size, int(total)))
	}
}

// Find auditLogs count histogram over time
// @Http: GET /histogram
// @QueryParam: from     | Timestamp           | start of time range filter
// @QueryParam: to       | Timestamp           | end of time range filter
// @QueryParam: userId   | string              | filter auditLogs by user id
// @QueryParam: action   | string              | filter auditLogs by action
// @QueryParam: itemType | string              | filter auditLogs by item type
// @QueryParam: itemId   | string              | filter auditLogs by item id
// @QueryParam: itemName | string              | filter auditLogs by item name
// @QueryParam: search   | string              | filter auditLogs by free text search
// @QueryParam: sort     | string              | sort results by field and direction: (e.g. time = sort by time asc, time- = sort by time desc)
// @QueryParam: page     | int                 | page number (for pagination)
// @QueryParam: size     | int                 | number of items per page (for pagination)
// @Return: EntityResponse<TimeSeries<float64>>
func (h *AuditLogsEndPoint) histogram(c *gin.Context) {
	// Get token data
	td := h.GetTokenData(c)
	if td == nil {
		return
	}

	p := s.AuditLogsFindParams{
		From:     h.GetParamAsTimestamp(c, "from", -1000*60*60*24*30),
		To:       h.GetParamAsTimestamp(c, "to", -1),
		UserId:   h.GetParamAsString(c, "userId", ""),
		Action:   h.GetParamAsString(c, "action", ""),
		ItemType: h.GetParamAsString(c, "itemType", ""),
		ItemId:   h.GetParamAsString(c, "itemId", ""),
		ItemName: h.GetParamAsString(c, "itemName", ""),
		Search:   h.GetParamAsString(c, "search", ""),
		Sort:     h.GetParamAsString(c, "sort", ""),
		Page:     h.GetParamAsInt(c, "page", 1),
		Size:     h.GetParamAsInt(c, "size", 100),
	}

	if result, err := h.service.Histogram(p); err != nil {
		c.JSON(http.StatusInternalServerError, rest.NewErrorResponse(err))
	} else {
		c.JSON(http.StatusOK, rest.NewEntityResponse(result))
	}
}

// endregion
