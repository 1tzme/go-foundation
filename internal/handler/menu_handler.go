package handler

import (
	"hot-coffee/internal/service"
	"hot-coffee/pkg/logger"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// TODO: Implement MenuHandler struct
type MenuHandler struct {
	menuService service.MenuServiceInterface
	logger      *logger.Logger
}

// TODO: Implement constructor with logger injection
func NewMenuHandler(menuService service.MenuServiceInterface, logger *logger.Logger) *MenuHandler {
	return &MenuHandler{
		menuService: menuService,
		logger:      logger.WithComponent("menu_handler"),
	}
}

// TODO: Implement GetAllMenuItems HTTP handler - GET /api/v1/menu
// - Call menu service to get all items
// - Return 200 OK with menu items list or 500 on error
// - Log HTTP request/response
func (h *MenuHandler) GetAllMenuItems(w http.ResponseWriter, r *http.Request) {
	reqCtx := &logger.RequestContext{
		Method:     r.Method,
		Path:       r.URL.Path,
		RemoteAddr: r.RemoteAddr,
		StartTime:  time.Now(),
	}
	h.logger.LogRequest(reqCtx)

	_, err := h.menuService.GetAllMenuItems()
	if err != nil {
		// logger err
		// writeErrorResponse()
		reqCtx.StatusCode = http.StatusInternalServerError
		h.logger.LogResponse(reqCtx)
		return
	}

	// writeJSONResponse()
	reqCtx.StatusCode = http.StatusOK
	h.logger.LogResponse(reqCtx)
}

// TODO: Implement CreateMenuItem HTTP handler - POST /api/v1/menu
// - Parse JSON request body
// - Validate request format
// - Call menu service to create item
// - Return 201 Created with item data or 400/500 on error
// - Log HTTP request/response
func (h *MenuHandler) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	reqCtx := &logger.RequestContext{
		Method:     r.Method,
		Path:       r.URL.Path,
		RemoteAddr: r.RemoteAddr,
		StartTime:  time.Now(),
	}
	h.logger.LogRequest(reqCtx)

	var createdReq service.CreateMenuItemRequest
	err := parseRequestBody(r, &createdReq)
	if err != nil {
		// logger warn
		// writeErrorResponse()
		reqCtx.StatusCode = http.StatusBadRequest
		h.logger.LogRequest(reqCtx)
		return
	}

	newID := generateMenuItemID(createdReq.Name)
	_, err = h.menuService.CreateMenuItem(newID, createdReq)
	if err != nil {
		// logger warn
		// writeErrorResponse()
		reqCtx.StatusCode = http.StatusBadRequest
		h.logger.LogRequest(reqCtx)
		return
	}

	// writeJSONResponse()
	reqCtx.StatusCode = http.StatusCreated
	h.logger.LogRequest(reqCtx)
}

// TODO: Implement UpdateMenuItem HTTP handler - PUT /api/v1/menu/{id}
// - Extract item ID from URL path
// - Parse JSON request body
// - Validate request format
// - Call menu service to update item
// - Return 200 OK or appropriate error status
// - Log HTTP request/response
func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	reqCtx := &logger.RequestContext{
		Method:     r.Method,
		Path:       r.URL.Path,
		RemoteAddr: r.RemoteAddr,
		StartTime:  time.Now(),
	}
	h.logger.LogRequest(reqCtx)

	id := extractIDFromPath(r)

	updateReq := service.UpdateMenuItemRequest{}
	err := parseRequestBody(r, &updateReq)
	if err != nil {
		// logger warn
		// writeErrorResponse()
		reqCtx.StatusCode = http.StatusBadRequest
		h.logger.LogRequest(reqCtx)
		return
	}

	err = h.menuService.UpdateMenuItem(id, updateReq)
	if err != nil {
		// logger warn
		if strings.Contains(err.Error(), "not found") {
			// writeErrorResponse
			reqCtx.StatusCode = http.StatusNotFound
		} else {
			// writeErrorResponse
			reqCtx.StatusCode = http.StatusBadRequest
		}
		h.logger.LogResponse(reqCtx)
		return
	}

	// writeJSONResponse()
	reqCtx.StatusCode = http.StatusOK
	h.logger.LogRequest(reqCtx)
}

// TODO: Implement DeleteMenuItem HTTP handler - DELETE /api/v1/menu/{id}
// - Extract item ID from URL path
// - Validate ID format
// - Call menu service to delete item
// - Return 204 No Content, 404 if not found, or 500 on error
// - Log HTTP request/response
func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	reqCtx := &logger.RequestContext{
		Method:     r.Method,
		Path:       r.URL.Path,
		RemoteAddr: r.RemoteAddr,
		StartTime:  time.Now(),
	}
	h.logger.LogRequest(reqCtx)

	id := extractIDFromPath(r)

	err := h.menuService.DeleteMenuItem(id)
	if err != nil {
		// logger warn
		// writeErrorResponse()
		reqCtx.StatusCode = http.StatusNotFound
		h.logger.LogResponse(reqCtx)
		return
	}

	// writeJSONResponse()
	reqCtx.StatusCode = http.StatusOK
	h.logger.LogResponse(reqCtx)
}

// TODO: Implement GetPopularItems HTTP handler - GET /api/v1/menu/aggregations/popular
// - Call menu service for popular items aggregation
// - Return 200 OK with aggregation data or 500 on error
// - Log HTTP request/response
// func (h *MenuHandler) GetPopularItems(w http.ResponseWriter, r *http.Request)

// TODO: Implement private helper methods
// - writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) - Write JSON response
// - writeErrorResponse(w http.ResponseWriter, statusCode int, message string) - Write error response
// - parseRequestBody(r *http.Request, target interface{}) error - Parse JSON request body
// - extractIDFromPath(r *http.Request) string - Extract ID from URL path
// - validateMenuItemID(id string) error - Validate menu item ID format

func generateMenuItemID(name string) string {
	cleaned := strings.ToLower(name)
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	cleaned = reg.ReplaceAllString(cleaned, "_")
	cleaned = strings.Trim(cleaned, "_")

	if cleaned == "" {
		cleaned = "ingredient"
	}

	return cleaned
}
