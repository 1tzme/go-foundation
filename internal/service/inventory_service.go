package service

import (
	"fmt"
	"hot-coffee/models"
	"hot-coffee/pkg/logger"
)

// Temporary placeholder for InventoryServiceInterface
type InventoryRepositoryInterface interface{}

type UpdateInventoryItemRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Quantity     int    `json:"quantity"`
	MinThreshold int    `json:"min_threshold"`
	Unit         string `json:"unit"`
}

type UpdateQuantityRequest struct {
	Quantity int `json:"quantity"`
}

type InventoryServiceInterface interface {
	GetAllInventoryItems() ([]*models.InventoryItem, error)
	UpdateInventoryItem(id string, req UpdateInventoryItemRequest) error
	GetLowStockItems() ([]*models.InventoryItem, error)
	UpdateQuantity(id string, req UpdateQuantityRequest) error
	// GetInventoryValue() (*models.InventoryValueAggregation, error)
}

type InventoryService struct {
	inventoryRepo InventoryRepositoryInterface
	logger        *logger.Logger
}

// NewInventoryService creates a new instance of InventoryService
func NewInventoryService(inventoryRepo InventoryRepositoryInterface, logger *logger.Logger) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
		logger:        logger.WithComponent("inventory_service"),
	}
}

// GetAllInventoryItems returns all inventory items (placeholder implementation)
func (s *InventoryService) GetAllInventoryItems() ([]*models.InventoryItem, error) {
	s.logger.Info("Fetching all inventory items")

	// Placeholder: return static data until repository is implemented
	items := []*models.InventoryItem{
		{
			ID:           "1",
			Name:         "Coffee Beans",
			Description:  "Premium Arabica beans",
			Quantity:     100,
			MinThreshold: 10,
			Unit:         "kg",
		},
		{
			ID:           "2",
			Name:         "Milk",
			Description:  "Whole milk",
			Quantity:     50,
			MinThreshold: 5,
			Unit:         "liters",
		},
	}
	s.logger.Info("Fetched inventory items", "count", len(items))
	return items, nil
}

// UpdateInventoryItem updates an existing inventory item
func (s *InventoryService) UpdateInventoryItem(id string, req UpdateInventoryItemRequest) error {
	s.logger.Info("Updating inventory item", "id", id, "name", req.Name)

	// Example validation: name and quantity must be present
	if req.Name == "" {
		s.logger.Warn("Update failed: name is required", "id", id)
		return fmt.Errorf("name is required")
	}
	if req.Quantity < 0 {
		s.logger.Warn("Update failed: quantity must be non-negative", "id", id, "quantity", req.Quantity)
		return fmt.Errorf("quantity must be non-negative")
	}

	// TODO: Add repository update logic here

	s.logger.Info("Inventory item updated (placeholder)", "id", id)
	return nil
}

// GetLowStockItems returns items that are below the minimum threshold
func (s *InventoryService) GetLowStockItems() ([]*models.InventoryItem, error) {
	s.logger.Info("Fetching low stock items")

	// Placeholder: return static data until repository is implemented
	lowStockItems := []*models.InventoryItem{
		{
			ID:          "2",
			Name:        "Milk",
			Description: "Whole milk",
			Quantity:    5,
			Unit:        "liters",
		},
	}
	s.logger.Info("Fetched low stock items", "count", len(lowStockItems))
	return lowStockItems, nil
}

// UpdateQuantity updates the quantity of an inventory item
func (s *InventoryService) UpdateQuantity(id string, req UpdateQuantityRequest) error {
	s.logger.Info("Updating inventory item quantity", "id", id, "quantity", req.Quantity)

	// Example validation: quantity must be non-negative
	if req.Quantity < 0 {
		s.logger.Warn("Update failed: quantity must be non-negative", "id", id, "quantity", req.Quantity)
		return fmt.Errorf("quantity must be non-negative")
	}

	// Placeholder: add repository update logic here
	s.logger.Info("Inventory item quantity updated (placeholder)", "id", id, "quantity", req.Quantity)
	return nil
}

// TODO: Implement GetInventoryValue method - Calculate inventory value
// - Call repository for inventory value aggregation
// - Apply business logic for valuation
// - Log calculation event
// func (s *InventoryService) GetInventoryValue() (*models.InventoryValueAggregation, error)

// Private business logic methods

// validateInventoryItemData validates the data for an inventory item update
func validateInventoryItemData(req UpdateInventoryItemRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Quantity < 0 {
		return fmt.Errorf("quantity must be non-negative")
	}
	if req.MinThreshold < 0 {
		return fmt.Errorf("minimum threshold must be non-negative")
	}
	if req.Unit == "" {
		return fmt.Errorf("unit is required")
	}
	return nil
}

// checkLowStockCondition checks if an inventory item is below its minimum threshold
// func checkLowStockCondition(item *models.InventoryItem) bool {
// 	// Check if the item's quantity is below the minimum threshold
// 	return item.Quantity < item.MinThreshold
// }

// validateInventoryItemID validates the ID of an inventory item
func validateQuantity(quantity int) error {
	if quantity < 0 {
		return fmt.Errorf("quantity must be non-negative")
	}
	return nil
}

// logLowStockWarning logs a warning if an inventory item is below its minimum threshold
// func logLowStockWarning(item *models.InventoryItem) {
// 	if item.Quantity < item.MinThreshold {
// 		// Log a warning if the item's quantity is below the minimum threshold
// 		fmt.Printf("Warning: Low stock for item %s (ID: %s). Current quantity: %d, Minimum threshold: %d\n",
// 			item.Name, item.ID, item.Quantity, item.MinThreshold)
// 	}
// }

// validateThreshold checks if the minimum threshold is valid
func validateThreshold(threshold int) error {
	if threshold < 0 {
		return fmt.Errorf("minimum threshold must be non-negative")
	}
	return nil
}
