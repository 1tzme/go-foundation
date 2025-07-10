package service

import (
	"fmt"
	"hot-coffee/internal/repositories"
	"hot-coffee/models"
	"hot-coffee/pkg/logger"
)

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
	CreateInventoryItem(id string, req UpdateInventoryItemRequest) error
	GetInventoryItem(id string) (*models.InventoryItem, error)
	DeleteInventoryItem(id string) error
}

// CreateInventoryItem adds a new inventory item
func (s *InventoryService) CreateInventoryItem(id string, req UpdateInventoryItemRequest) error {
	s.logger.Info("Creating inventory item", "id", id, "name", req.Name)
	if err := validateInventoryItemData(req); err != nil {
		s.logger.Warn("Create failed: invalid data", "id", id, "error", err)
		return err
	}
	item := &models.InventoryItem{
		IngredientID: id,
		Name:         req.Name,
		Quantity:     float64(req.Quantity),
		MinThreshold: float64(req.MinThreshold),
		Unit:         req.Unit,
	}
	if err := s.inventoryRepo.Add(item); err != nil {
		s.logger.Error("Failed to add inventory item in repository", "id", id, "error", err)
		return err
	}
	s.logger.Info("Inventory item created", "id", id)
	return nil
}

// GetInventoryItem fetches a single inventory item by ID
func (s *InventoryService) GetInventoryItem(id string) (*models.InventoryItem, error) {
	s.logger.Info("Fetching inventory item by id", "id", id)
	item, err := s.inventoryRepo.GetByID(id)
	if err != nil {
		s.logger.Warn("Inventory item not found", "id", id, "error", err)
		return nil, err
	}
	s.logger.Info("Fetched inventory item", "id", id)
	return item, nil
}

// DeleteInventoryItem deletes an inventory item by ID
func (s *InventoryService) DeleteInventoryItem(id string) error {
	s.logger.Info("Deleting inventory item", "id", id)
	if err := s.inventoryRepo.Delete(id); err != nil {
		s.logger.Warn("Failed to delete inventory item", "id", id, "error", err)
		return err
	}
	s.logger.Info("Inventory item deleted", "id", id)
	return nil
}

type InventoryService struct {
	inventoryRepo repositories.InventoryRepositoryInterface
	logger        *logger.Logger
}

// NewInventoryService creates a new instance of InventoryService
func NewInventoryService(inventoryRepo repositories.InventoryRepositoryInterface, logger *logger.Logger) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
		logger:        logger.WithComponent("inventory_service"),
	}
}

// GetAllInventoryItems returns all inventory items (placeholder implementation)
func (s *InventoryService) GetAllInventoryItems() ([]*models.InventoryItem, error) {
	s.logger.Info("Fetching all inventory items from repository")
	items, err := s.inventoryRepo.GetAll()
	if err != nil {
		s.logger.Error("Failed to fetch inventory items from repository", "error", err)
		return nil, err
	}
	s.logger.Info("Fetched inventory items", "count", len(items))
	return items, nil
}

// UpdateInventoryItem updates an existing inventory item
func (s *InventoryService) UpdateInventoryItem(id string, req UpdateInventoryItemRequest) error {
	s.logger.Info("Updating inventory item", "id", id, "name", req.Name)

	// Validate input
	if err := validateInventoryItemData(req); err != nil {
		s.logger.Warn("Update failed: invalid data", "id", id, "error", err)
		return err
	}

	// Build item struct for update
	item := &models.InventoryItem{
		IngredientID: id,
		Name:         req.Name,
		Quantity:     float64(req.Quantity),
		MinThreshold: float64(req.MinThreshold),
		Unit:         req.Unit,
	}

	err := s.inventoryRepo.Update(id, item)
	if err != nil {
		s.logger.Error("Failed to update inventory item in repository", "id", id, "error", err)
		return err
	}

	s.logger.Info("Inventory item updated", "id", id)
	return nil
}

// GetLowStockItems returns items that are below the minimum threshold
func (s *InventoryService) GetLowStockItems() ([]*models.InventoryItem, error) {
	s.logger.Info("Fetching low stock items from repository")
	items, err := s.inventoryRepo.GetLowStockItems()
	if err != nil {
		s.logger.Error("Failed to fetch low stock items from repository", "error", err)
		return nil, err
	}
	s.logger.Info("Fetched low stock items", "count", len(items))
	return items, nil
}

// UpdateQuantity updates the quantity of an inventory item
func (s *InventoryService) UpdateQuantity(id string, req UpdateQuantityRequest) error {
	s.logger.Info("Updating inventory item quantity", "id", id, "quantity", req.Quantity)

	if err := validateQuantity(req.Quantity); err != nil {
		s.logger.Warn("Update failed: quantity must be non-negative", "id", id, "quantity", req.Quantity)
		return err
	}

	err := s.inventoryRepo.UpdateQuantity(id, req.Quantity)
	if err != nil {
		s.logger.Error("Failed to update inventory item quantity in repository", "id", id, "error", err)
		return err
	}

	s.logger.Info("Inventory item quantity updated", "id", id, "quantity", req.Quantity)
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
