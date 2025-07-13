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

	items, err := h.menuService.GetAllMenuItems()
	if err != nil {
		h.logger.Error("Failed to get all menu items")
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to fetch menu items")
		reqCtx.StatusCode = http.StatusInternalServerError
		h.logger.LogResponse(reqCtx)
		return
	}

	writeJSONResponse(w, http.StatusOK, items)
	reqCtx.StatusCode = http.StatusOK
	h.logger.LogResponse(reqCtx)
}

func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
	reqCtx := &logger.RequestContext{
		Method: r.Method,
		Path: r.URL.Path,
		RemoteAddr: r.RemoteAddr,
		StartTime: time.Now(),
	}
	h.logger.LogRequest(reqCtx)

	id := extractIDFromPath(r)
	item, err := h.menuService.GetMenuItem(id)
	if err != nil {
		h.logger.Warn("Menu item not found", "id", id, "error", err)
		writeErrorResponse(w, http.StatusNotFound, "Menu item not found")
		reqCtx.StatusCode = http.StatusNotFound
		h.logger.LogResponse(reqCtx)
		return
	}

	writeJSONResponse(w, http.StatusOK, item)
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
	if err := parseRequestBody(r, &createdReq); err != nil {
		h.logger.Warn("Invalid request body for create", "error", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		reqCtx.StatusCode = http.StatusBadRequest
		h.logger.LogResponse(reqCtx)
		return
	}

	newID := generateMenuItemID(createdReq.Name)
	item, err := h.menuService.CreateMenuItem(newID, createdReq)
	if err != nil {
		h.logger.Warn("Failed to create menu item", "error", err)
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		reqCtx.StatusCode = http.StatusBadRequest
		h.logger.LogResponse(reqCtx)
		return
	}

	writeJSONResponse(w, http.StatusCreated, map[string]interface{}{"id": newID, "message": "Menu item created", "item": item})
	reqCtx.StatusCode = http.StatusCreated
	h.logger.LogResponse(reqCtx)
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
	if err := parseRequestBody(r, &updateReq); err != nil {
		h.logger.Warn("Invalid request body", "error", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		reqCtx.StatusCode = http.StatusBadRequest
		h.logger.LogResponse(reqCtx)
		return
	}

	if err := h.menuService.UpdateMenuItem(id, updateReq); err != nil {
		h.logger.Warn("Failed to update menu item", "id", id, "error", err)
		if strings.Contains(err.Error(), "not found") {
			writeErrorResponse(w, http.StatusNotFound, err.Error())
			reqCtx.StatusCode = http.StatusNotFound
		} else {
			writeErrorResponse(w, http.StatusBadRequest, err.Error())
			reqCtx.StatusCode = http.StatusBadRequest
		}
		h.logger.LogResponse(reqCtx)
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{"id": id, "message": "Menu item updated"})
	reqCtx.StatusCode = http.StatusOK
	h.logger.LogResponse(reqCtx)
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

	if err := h.menuService.DeleteMenuItem(id); err != nil {
		h.logger.Warn("Failed to delete menu item", "id", id, "error", err)
		writeErrorResponse(w, http.StatusNotFound, "Menu item not found")
		reqCtx.StatusCode = http.StatusNotFound
		h.logger.LogResponse(reqCtx)
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{"id": id, "message": "Menu item deleted"})
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
