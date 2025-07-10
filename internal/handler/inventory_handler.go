package handler

import (
	"encoding/json"
	"fmt"
	"hot-coffee/pkg/logger"
	"net/http"
	"time"
)

// Temporary placeholder for InventoryServiceInterface
type InventoryServiceInterface interface{}

type InventoryHandler struct {
	inventoryService InventoryServiceInterface
	logger           *logger.Logger
}

func NewInventoryHandler(inventoryService InventoryServiceInterface, logger *logger.Logger) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
		logger:           logger.WithComponent("inventory_handler"),
	}
}

// TODO: Implement GetAllInventoryItems HTTP handler - GET /api/v1/inventory
// - Call inventory service to get all items
// - Return 200 OK with inventory items list or 500 on error
// - Log HTTP request/response
func (h *InventoryHandler) GetAllInventoryItems(w http.ResponseWriter, r *http.Request) {
	reqCtx := &logger.RequestContext{
		Method:     r.Method,
		Path:       r.URL.Path,
		RemoteAddr: r.RemoteAddr,
		StartTime:  time.Now(),
	}
	h.logger.LogRequest(reqCtx)

	// Placeholder: return static data until service is ready
	items := []map[string]interface{}{
		{"id": 1, "name": "Coffee Beans", "quantity": 100},
		{"id": 2, "name": "Milk", "quantity": 50},
	}
	writeJSONResponse(w, http.StatusOK, items)

	reqCtx.StatusCode = http.StatusOK
	reqCtx.ResponseSize = int64(len(items)) // For real use, calculate bytes written
	h.logger.LogResponse(reqCtx)
}

// GetInventory is a basic HTTP handler for GET /api/v1/inventory
func (h *InventoryHandler) GetInventory(w http.ResponseWriter, r *http.Request) {
	reqCtx := &logger.RequestContext{
		Method:     r.Method,
		Path:       r.URL.Path,
		RemoteAddr: r.RemoteAddr,
		StartTime:  time.Now(),
	}
	h.logger.LogRequest(reqCtx)

	// Placeholder: return static data until service is ready
	items := []map[string]interface{}{
		{"id": 1, "name": "Coffee Beans", "quantity": 100},
		{"id": 2, "name": "Milk", "quantity": 50},
	}
	writeJSONResponse(w, http.StatusOK, items)

	reqCtx.StatusCode = http.StatusOK
	reqCtx.ResponseSize = int64(len(items)) // For real use, calculate bytes written
	h.logger.LogResponse(reqCtx)
}

// UpdateInventoryItem handles PUT /api/v1/inventory/{id}
func (h *InventoryHandler) UpdateInventoryItem(w http.ResponseWriter, r *http.Request) {
	reqCtx := &logger.RequestContext{
		Method:     r.Method,
		Path:       r.URL.Path,
		RemoteAddr: r.RemoteAddr,
		StartTime:  time.Now(),
	}
	h.logger.LogRequest(reqCtx)

	id := extractIDFromPath(r)
	if err := validateInventoryItemID(id); err != nil {
		h.logger.Warn("Invalid inventory item ID", "id", id, "error", err)
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		reqCtx.StatusCode = http.StatusBadRequest
		h.logger.LogResponse(reqCtx)
		return
	}

	// Placeholder for request body struct
	var updateReq map[string]interface{}
	if err := parseRequestBody(r, &updateReq); err != nil {
		h.logger.Warn("Invalid request body", "error", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		reqCtx.StatusCode = http.StatusBadRequest
		h.logger.LogResponse(reqCtx)
		return
	}

	// TODO: Call inventoryService.UpdateItem(id, updateReq) and handle result
	// For now, just echo back the update
	resp := map[string]interface{}{
		"id":   id,
		"data": updateReq,
	}
	writeJSONResponse(w, http.StatusOK, resp)
	reqCtx.StatusCode = http.StatusOK
	h.logger.LogResponse(reqCtx)
}

// GetLowStockItems handles GET /api/v1/inventory/low-stock
func (h *InventoryHandler) GetLowStockItems(w http.ResponseWriter, r *http.Request) {
	reqCtx := &logger.RequestContext{
		Method:     r.Method,
		Path:       r.URL.Path,
		RemoteAddr: r.RemoteAddr,
		StartTime:  time.Now(),
	}
	h.logger.LogRequest(reqCtx)

	// Placeholder: static low stock items
	lowStock := []map[string]interface{}{
		{"id": 2, "name": "Milk", "quantity": 5},
	}

	if len(lowStock) > 0 {
		h.logger.Warn("Low stock items found", "count", len(lowStock))
	}

	writeJSONResponse(w, http.StatusOK, lowStock)
	reqCtx.StatusCode = http.StatusOK
	reqCtx.ResponseSize = int64(len(lowStock))
	h.logger.LogResponse(reqCtx)
}

// UpdateQuantity handles PATCH /api/v1/inventory/{id}/quantity
func (h *InventoryHandler) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	reqCtx := &logger.RequestContext{
		Method:     r.Method,
		Path:       r.URL.Path,
		RemoteAddr: r.RemoteAddr,
		StartTime:  time.Now(),
	}
	h.logger.LogRequest(reqCtx)

	id := extractIDFromPath(r)
	if err := validateInventoryItemID(id); err != nil {
		h.logger.Warn("Invalid inventory item ID", "id", id, "error", err)
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		reqCtx.StatusCode = http.StatusBadRequest
		h.logger.LogResponse(reqCtx)
		return
	}

	// Parse quantity from request body
	var reqBody struct {
		Quantity int `json:"quantity"`
	}
	if err := parseRequestBody(r, &reqBody); err != nil {
		h.logger.Warn("Invalid request body for quantity update", "error", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		reqCtx.StatusCode = http.StatusBadRequest
		h.logger.LogResponse(reqCtx)
		return
	}
	if reqBody.Quantity < 0 {
		h.logger.Warn("Invalid quantity value", "quantity", reqBody.Quantity)
		writeErrorResponse(w, http.StatusBadRequest, "Quantity must be non-negative")
		reqCtx.StatusCode = http.StatusBadRequest
		h.logger.LogResponse(reqCtx)
		return
	}

	// TODO: Call inventoryService.UpdateQuantity(id, reqBody.Quantity) and handle result
	// For now, just echo back the update
	resp := map[string]interface{}{
		"id":       id,
		"quantity": reqBody.Quantity,
	}
	writeJSONResponse(w, http.StatusOK, resp)
	reqCtx.StatusCode = http.StatusOK
	h.logger.LogResponse(reqCtx)
}

// GetInventoryValue handles GET /api/v1/inventory/value
func (h *InventoryHandler) GetInventoryValue(w http.ResponseWriter, r *http.Request) {
	reqCtx := &logger.RequestContext{
		Method:     r.Method,
		Path:       r.URL.Path,
		RemoteAddr: r.RemoteAddr,
		StartTime:  time.Now(),
	}
	h.logger.LogRequest(reqCtx)

	// Placeholder: static inventory value
	value := map[string]interface{}{
		"total_value": 1234.56,
		"currency":    "USD",
	}
	writeJSONResponse(w, http.StatusOK, value)
	reqCtx.StatusCode = http.StatusOK
	h.logger.LogResponse(reqCtx)
}

// Private helper methods

// writeJSONResponse - writes JSON response with given status code and data
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
		}
	}
}

// writeErrorResponse - writes an error response with given status code and message
func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := map[string]string{"error": message}
	_ = json.NewEncoder(w).Encode(resp)
}

// parseRequestBody - parses JSON request body into the target struct
func parseRequestBody(r *http.Request, target interface{}) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(target)
}

// extractIDFromPath - extracts ID from URL path (expects /api/v1/inventory/{id} or similar)
func extractIDFromPath(r *http.Request) string {
	parts := splitPath(r.URL.Path)
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// splitPath - splits a URL path into segments, removing empty segments
func splitPath(path string) []string {
	var segments []string
	for _, p := range split(path, '/') {
		if p != "" {
			segments = append(segments, p)
		}
	}
	return segments
}

// split - splits a string by a separator rune
func split(s string, sep rune) []string {
	var res []string
	last := 0
	for i, c := range s {
		if c == sep {
			res = append(res, s[last:i])
			last = i + 1
		}
	}
	res = append(res, s[last:])
	return res
}

// validateInventoryItemID - Validate inventory item ID format
func validateInventoryItemID(id string) error {
	if id == "" {
		return fmt.Errorf("item ID is required")
	}
	for _, c := range id {
		if c < '0' || c > '9' {
			return fmt.Errorf("invalid item ID: must be numeric")
		}
	}
	return nil
}
