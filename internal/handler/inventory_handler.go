package handler

import (
	"encoding/json"
	"fmt"
	"hot-coffee/internal/service"
	"hot-coffee/pkg/logger"
	"net/http"
	"time"
)

type InventoryHandler struct {
	inventoryService service.InventoryServiceInterface
	logger           *logger.Logger
}

// CreateInventoryItem handles POST /api/v1/inventory
func (h *InventoryHandler) CreateInventoryItem(w http.ResponseWriter, r *http.Request) {
	reqCtx := &logger.RequestContext{
		Method:     r.Method,
		Path:       r.URL.Path,
		RemoteAddr: r.RemoteAddr,
		StartTime:  time.Now(),
	}
	h.logger.LogRequest(reqCtx)

	var createReq service.UpdateInventoryItemRequest
	if err := parseRequestBody(r, &createReq); err != nil {
		h.logger.Warn("Invalid request body for create", "error", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		reqCtx.StatusCode = http.StatusBadRequest
		h.logger.LogResponse(reqCtx)
		return
	}

	// For now, generate a new ID as a timestamp string (replace with better ID logic as needed)
	newID := fmt.Sprintf("%d", time.Now().UnixNano())
	err := h.inventoryService.UpdateInventoryItem(newID, createReq)
	if err != nil {
		h.logger.Warn("Failed to create inventory item", "error", err)
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		reqCtx.StatusCode = http.StatusBadRequest
		h.logger.LogResponse(reqCtx)
		return
	}

	writeJSONResponse(w, http.StatusCreated, map[string]interface{}{"id": newID, "message": "Inventory item created"})
	reqCtx.StatusCode = http.StatusCreated
	h.logger.LogResponse(reqCtx)
}

// GetInventoryItem handles GET /api/v1/inventory/{id}
func (h *InventoryHandler) GetInventoryItem(w http.ResponseWriter, r *http.Request) {
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

	items, err := h.inventoryService.GetAllInventoryItems()
	if err != nil {
		h.logger.Error("Failed to fetch inventory items for get by id", "error", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to fetch inventory items")
		reqCtx.StatusCode = http.StatusInternalServerError
		h.logger.LogResponse(reqCtx)
		return
	}
	var found interface{} = nil
	for _, item := range items {
		if item.IngredientID == id {
			found = item
			break
		}
	}
	if found == nil {
		h.logger.Warn("Inventory item not found", "id", id)
		writeErrorResponse(w, http.StatusNotFound, "Inventory item not found")
		reqCtx.StatusCode = http.StatusNotFound
		h.logger.LogResponse(reqCtx)
		return
	}
	writeJSONResponse(w, http.StatusOK, found)
	reqCtx.StatusCode = http.StatusOK
	h.logger.LogResponse(reqCtx)
}

// DeleteInventoryItem handles DELETE /api/v1/inventory/{id}
func (h *InventoryHandler) DeleteInventoryItem(w http.ResponseWriter, r *http.Request) {
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

	// For now, simulate delete by removing from in-memory list (should call service/repo in real impl)
	items, err := h.inventoryService.GetAllInventoryItems()
	if err != nil {
		h.logger.Error("Failed to fetch inventory items for delete", "error", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to fetch inventory items")
		reqCtx.StatusCode = http.StatusInternalServerError
		h.logger.LogResponse(reqCtx)
		return
	}
	found := false
	for _, item := range items {
		if item.IngredientID == id {
			found = true
			break
		}
	}
	if !found {
		h.logger.Warn("Inventory item not found for delete", "id", id)
		writeErrorResponse(w, http.StatusNotFound, "Inventory item not found")
		reqCtx.StatusCode = http.StatusNotFound
		h.logger.LogResponse(reqCtx)
		return
	}
	// TODO: Call service/repo to actually delete item
	writeJSONResponse(w, http.StatusOK, map[string]interface{}{"id": id, "message": "Inventory item deleted (not really, stub)"})
	reqCtx.StatusCode = http.StatusOK
	h.logger.LogResponse(reqCtx)
}

// NewInventoryHandler creates a new InventoryHandler with the given inventory service and logger
func NewInventoryHandler(inventoryService service.InventoryServiceInterface, logger *logger.Logger) *InventoryHandler {
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

	items, err := h.inventoryService.GetAllInventoryItems()
	if err != nil {
		h.logger.Error("Failed to get all inventory items", "error", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to fetch inventory items")
		reqCtx.StatusCode = http.StatusInternalServerError
		h.logger.LogResponse(reqCtx)
		return
	}

	writeJSONResponse(w, http.StatusOK, items)
	reqCtx.StatusCode = http.StatusOK
	// Optionally, calculate bytes written for ResponseSize
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

	items, err := h.inventoryService.GetAllInventoryItems()
	if err != nil {
		h.logger.Error("Failed to get inventory", "error", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to fetch inventory items")
		reqCtx.StatusCode = http.StatusInternalServerError
		h.logger.LogResponse(reqCtx)
		return
	}

	writeJSONResponse(w, http.StatusOK, items)
	reqCtx.StatusCode = http.StatusOK
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

	var updateReq service.UpdateInventoryItemRequest
	if err := parseRequestBody(r, &updateReq); err != nil {
		h.logger.Warn("Invalid request body", "error", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		reqCtx.StatusCode = http.StatusBadRequest
		h.logger.LogResponse(reqCtx)
		return
	}

	err := h.inventoryService.UpdateInventoryItem(id, updateReq)
	if err != nil {
		h.logger.Warn("Failed to update inventory item", "id", id, "error", err)
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		reqCtx.StatusCode = http.StatusBadRequest
		h.logger.LogResponse(reqCtx)
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{"id": id, "message": "Inventory item updated"})
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

	items, err := h.inventoryService.GetLowStockItems()
	if err != nil {
		h.logger.Error("Failed to get low stock items", "error", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to fetch low stock items")
		reqCtx.StatusCode = http.StatusInternalServerError
		h.logger.LogResponse(reqCtx)
		return
	}

	if len(items) > 0 {
		h.logger.Warn("Low stock items found", "count", len(items))
	}

	writeJSONResponse(w, http.StatusOK, items)
	reqCtx.StatusCode = http.StatusOK
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

	var reqBody service.UpdateQuantityRequest
	if err := parseRequestBody(r, &reqBody); err != nil {
		h.logger.Warn("Invalid request body for quantity update", "error", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		reqCtx.StatusCode = http.StatusBadRequest
		h.logger.LogResponse(reqCtx)
		return
	}

	err := h.inventoryService.UpdateQuantity(id, reqBody)
	if err != nil {
		h.logger.Warn("Failed to update quantity", "id", id, "error", err)
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		reqCtx.StatusCode = http.StatusBadRequest
		h.logger.LogResponse(reqCtx)
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{"id": id, "quantity": reqBody.Quantity, "message": "Quantity updated"})
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
	return nil
}
