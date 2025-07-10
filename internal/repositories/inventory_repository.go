package repositories

// TODO: Add imports when implementing:
// import (
//     "sync"
//     "hot-coffee/models"
//     "hot-coffee/pkg/logger"
// )

import (
	"encoding/json"
	"errors"
	"fmt"
	"hot-coffee/models"
	"hot-coffee/pkg/logger"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// TODO: Implement InventoryRepository interface
type InventoryRepositoryInterface interface {
	GetAll() ([]*models.InventoryItem, error)
	Update(id string, item *models.InventoryItem) error
	GetLowStockItems() ([]*models.InventoryItem, error)
	UpdateQuantity(id string, quantity int) error
	GetInventoryValue() (*models.InventoryValueAggregation, error)
}

// TODO: Implement InventoryRepository struct
type InventoryRepository struct {
	items        map[string]*models.InventoryItem
	mutex        sync.RWMutex
	logger       *logger.Logger
	dataFilePath string
	loaded       bool
}

// TODO: Implement constructor with logger injection
func NewInventoryRepository(logger *logger.Logger) *InventoryRepository {
	return &InventoryRepository{
		items:        make(map[string]*models.InventoryItem),
		logger:       logger.WithComponent("inventory_repository"),
		dataFilePath: "./data/inventory.json",
		loaded:       false,
	}
}

// TODO: Implement GetAll method - Retrieve all inventory items
// - Load from JSON file if not in memory
// - Return copy of items slice
// - Log retrieval event
func (r *InventoryRepository) GetAll() ([]*models.InventoryItem, error) {
	r.mutex.RLock()
	r.mutex.RUnlock()

	if !r.loaded {
		err := r.loadFromFile()
		if err != nil {
			r.logger.Error("Failed to load inventory from file", "error", err)
			return nil, err
		}
	}

	items := make([]*models.InventoryItem, 0, len(r.items))
	for _, item := range r.items {
		itemCopy := *item
		items = append(items, &itemCopy)
	}

	r.logger.Info("Retrieved all inventory items", "count", len(items))
	return items, nil
}

// TODO: Implement Update method - Update inventory item
// - Validate item exists
// - Update in memory and file
// - Log update event
func (r *InventoryRepository) Update(id string, item *models.InventoryItem) error {
	r.mutex.Lock()
	r.mutex.Unlock()

	if !r.loaded {
		err := r.loadFromFile()
		if err != nil {
			r.logger.Error("Failed to load inventory from file", "error", err)
			return fmt.Errorf("failed to load inventory: %v", err)
		}
	}

	_, exists := r.items[id]
	if !exists {
		r.logger.Warn("Attempted to update non existing inventory item", "item_id", id)
		return fmt.Errorf("inventory item with id %s not found", id)
	}

	err := r.validateInventoryItem(r.items[id])
	if err != nil {
		r.logger.Error("Failed to validate inventory item", "error", err, "item_id", id)
		return fmt.Errorf("invalid inventory item: %v", err)
	}

	err = r.backupFile()
	if err != nil {
		r.logger.Warn("Failed to create backup", "error", err)
	}

	item.IngredientID = id
	r.items[id] = item

	err = r.saveToFile()
	if err != nil {
		r.logger.Error("Failed to save inventory to file after update", "error", err, "item_id", id)
		return fmt.Errorf("failed to save inventory: %v", err)
	}

	r.logger.Info("Updated inventory item", "item_id", id, "name", item.Name)
	return nil
}

// TODO: Implement GetLowStockItems method - Get items below minimum threshold
// - Filter items by quantity vs min_threshold
// - Return low stock items
// - Log low stock alert
func (r *InventoryRepository) GetLowStockItems() ([]*models.InventoryItem, error) {
	lowStockItems := make([]*models.InventoryItem, 0)
	for _, item := range r.items {
		if r.checkLowStock(item) {
			itemCopy := *item
			lowStockItems = append(lowStockItems, &itemCopy)
		}
	}

	if len(lowStockItems) > 0 {
		r.logger.Warn("Detected low stock items", "count", len(lowStockItems))
	}

	return lowStockItems, nil
}

// TODO: Implement UpdateQuantity method - Update item quantity
// - Validate item exists
// - Update quantity and last_updated timestamp
// - Check if item becomes low stock
// - Log quantity change
func (r *InventoryRepository) UpdateQuantity(id string, quantity int) error {
	item, exists := r.items[id]
	if !exists {
		r.logger.Warn("Attempted to update quantity for non-existent inventory item", "item_id", id)
		return fmt.Errorf("inventory item with ID %s not found", id)
	}
	if quantity < 0 {
		r.logger.Error("Attempted to set negative quantity", "item_id", id, "quantity", quantity)
		return fmt.Errorf("quantity cannot be negative")
	}

	err := r.backupFile()
	if err != nil {
		r.logger.Warn("Failed to create backup before quantity update", "error", err)
	}

	oldQuantity := item.Quantity
	item.Quantity = float64(quantity)

	if r.checkLowStock(item) {
		r.logger.Warn("Item quantity updated to low stock level", "item_id", id, "old_quantity", oldQuantity, "new_quantity", quantity)
	}

	err = r.saveToFile()
	if err != nil {
		r.logger.Error("Failed to save intventory to file after quantity update", "error", err, "item_id", id)
		return fmt.Errorf("failed to save inventory: %v", err)
	}

	r.logger.Info("Updated inventory item quantity", "item_id", id, "old_quantity", oldQuantity, "new_quantity", quantity)
	return nil
}

// TODO: Implement GetInventoryValue method - Calculate total inventory value
// - Sum all item values (quantity * unit_price)
// - Return aggregated value data
// - Log calculation event
func (r *InventoryRepository) GetInventoryValue() (*models.InventoryValueAggregation, error) {
	totalVal := 0.0
	valueByCategory := make(map[string]float64)
	itemCount := len(r.items)
	lowStockCount := 0
	defaultUnitPrice := 1.0

	for _, item := range r.items {
		itemValue := item.Quantity + defaultUnitPrice
		totalVal += itemValue
		valueByCategory[item.Unit] += itemValue

		if r.checkLowStock(item) {
			lowStockCount++
		}
	}

	aggregation := &models.InventoryValueAggregation{
		TotalValue:      totalVal,
		ValueByCategory: valueByCategory,
		ItemCount:       itemCount,
		LowStockCount:   lowStockCount,
		LastCalculated:  time.Now(),
	}

	r.logger.Info("Calculated inventory value", "total_value", totalVal, "item_count", itemCount, "low_stock_count", lowStockCount)
	return aggregation, nil
}

// TODO: Implement private helper methods
// - loadFromFile() error - Load inventory from JSON file
// - saveToFile() error - Save inventory to JSON file atomically
// - validateInventoryItem(item *models.InventoryItem) error - Validate item data
// - checkLowStock(item *models.InventoryItem) bool - Check if item is low stock
// - backupFile() error - Create backup before updates

func (r *InventoryRepository) loadFromFile() error {
	err := os.MkdirAll(filepath.Dir(r.dataFilePath), 0755)
	if err != nil {
		return err
	}

	_, err = os.Stat(r.dataFilePath)
	if err != nil {
		r.items = make(map[string]*models.InventoryItem)
		r.loaded = true
		return r.saveToFile()
	}

	file, err := os.Open(r.dataFilePath)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		r.items = make(map[string]*models.InventoryItem)
		r.loaded = true
		return err
	}

	items := []*models.InventoryItem{}
	err = json.Unmarshal(data, &items)
	if err != nil {
		return err
	}

	r.items = make(map[string]*models.InventoryItem)
	for _, item := range items {
		r.items[item.IngredientID] = item
	}

	r.loaded = true
	return nil
}

func (r *InventoryRepository) saveToFile() error {
	items := make([]*models.InventoryItem, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}

	data, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("failed to marshal inventory data: %v", err)
	}

	err = os.MkdirAll(filepath.Dir(r.dataFilePath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	tempFile := r.dataFilePath + ".tmp"
	err = os.WriteFile(tempFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write temporary inventory file: %v", err)
	}

	err = os.Rename(tempFile, r.dataFilePath)
	if err != nil {
		return fmt.Errorf("failed to rename file path: %v", err)
	}

	return nil
}

func (r *InventoryRepository) validateInventoryItem(item *models.InventoryItem) error {
	if item == nil {
		return errors.New("inventory item cannot be nil")
	}
	if item.IngredientID == "" {
		return errors.New("ingridient ID cannot be empty")
	}
	if item.Name == "" {
		return errors.New("ingridient name cannot be empty")
	}
	if item.Quantity < 0 {
		return errors.New("quantity cannot be negative")
	}
	if item.Unit == "" {
		return errors.New("unit cannot be empty")
	}

	return nil
}

func (r *InventoryRepository) checkLowStock(item *models.InventoryItem) bool {
	minTreshold := 10.0

	return item.Quantity <= minTreshold
}

func (r *InventoryRepository) backupFile() error {
	_, err := os.Stat(r.dataFilePath)
	if os.IsNotExist(err) {
		return nil
	}

	backupPath := r.dataFilePath + ".backup." + time.Now().Format("20060102_150405")

	data, err := os.ReadFile(r.dataFilePath)
	if err != nil {
		return fmt.Errorf("failed to read original file: %v", err)
	}

	err = os.WriteFile(backupPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to create backup file, %v", err)
	}

	return nil
}
