package repositories

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

// TODO: Add imports when implementing:
// import (
//     "sync"
//     "hot-coffee/models"
//     "hot-coffee/pkg/logger"
// )

// TODO: Implement MenuRepository interface
type MenuRepositoryInterface interface {
	GetAll() ([]*models.MenuItem, error)
	Create(item *models.MenuItem) error
	Update(id string, item *models.MenuItem) error
	Delete(id string) error
	GetByID(id string) (*models.MenuItem, error)
}

// TODO: Implement MenuRepository struct
type MenuRepository struct {
	items        map[string]*models.MenuItem
	mutex        sync.RWMutex
	logger       *logger.Logger
	dataFilePath string
	loaded       bool
}

// TODO: Implement constructor with logger injection
func NewMenuRepository(logger *logger.Logger) *MenuRepository {
	return &MenuRepository{
		items:        make(map[string]*models.MenuItem),
		logger:       logger.WithComponent("menu_repository"),
		dataFilePath: "./data/menu_items.json",
		loaded:       false,
	}
}

// TODO: Implement GetAll method - Retrieve all menu items
// - Load from JSON file if not in memory
// - Return copy of items slice
// - Log retrieval event
func (r *MenuRepository) GetAll() ([]*models.MenuItem, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.loaded {
		err := r.loadFromFile()
		if err != nil {
			// logger err
			return nil, err
		}
	}

	items := make([]*models.MenuItem, 0, len(r.items))
	for _, item := range r.items {
		itemCopy := *item
		items = append(items, &itemCopy)
	}

	// logger info
	return items, nil
}

// TODO: Implement Create method - Create a new menu item
// - Generate unique item ID
// - Validate item data
// - Save to memory map
// - Persist to JSON file atomically
// - Log creation event
func (r *MenuRepository) Create(item *models.MenuItem) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.loaded {
		err := r.loadFromFile()
		if err != nil {
			// logger err
			return err
		}
	}

	_, exists := r.items[item.ID]
	if exists {
		// logger warn
		return fmt.Errorf("menu item with ID %s already exists", item.ID)
	}

	err := r.validateMenuItem(item)
	if err != nil {
		// logger err
		return err
	}

	r.items[item.ID] = item

	err = r.saveToFile()
	if err != nil {
		// logger err
		return err
	}

	// logger info
	return nil
}

// TODO: Implement Update method - Update existing menu item
// - Validate item exists
// - Update in memory and file
// - Log update event
func (r *MenuRepository) Update(id string, item *models.MenuItem) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.loaded {
		err := r.loadFromFile()
		if err != nil {
			// logger err
			return err
		}
	}

	_, exists := r.items[id]
	if !exists {
		// logger warn
		return fmt.Errorf("menu item with id %s not found", id)
	}

	err := r.validateMenuItem(item)
	if err != nil {
		// logger err
		return err
	}
	err = r.backupFile()
	if err != nil {
		// logger warn
	}
	
	item.ID = id
	r.items[id] = item

	err = r.saveToFile()
	if err != nil {
		// logger err
		return err
	}

	// logger info
	return nil
}

// TODO: Implement Delete method - Delete menu item
// - Validate item exists
// - Remove from memory and file
// - Log deletion event
func (r *MenuRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.loaded {
		err := r.loadFromFile()
		if err != nil {
			// logger err
			return err
		}
	}

	_, exists := r.items[id]
	if !exists {
		// logger warn
		return fmt.Errorf("menu item with id %s not found", id)
	}
	err := r.backupFile()
	if err != nil {
		return err
	}

	delete(r.items, id)

	err = r.saveToFile()
	if err != nil {
		// logger err
		return err
	}

	// logger info
	return nil
}

// TODO: Implement GetPopularItems method - Get popular menu items aggregation
// - Analyze order history
// - Count item frequencies
// - Return sorted popular items
// func (r *MenuRepository) GetPopularItems() ([]*models.PopularItemAggregation, error)

// TODO: Implement private helper methods
// - loadFromFile() error - Load menu items from JSON file
// - saveToFile() error - Save menu items to JSON file atomically
// - generateItemID() string - Generate unique item ID
// - validateMenuItem(item *models.MenuItem) error - Validate item data
// - backupFile() error - Create backup before updates

func (r *MenuRepository) loadFromFile() error {
	err := os.MkdirAll(filepath.Dir(r.dataFilePath), 0755)
	if err != nil {
		return err
	}

	_, err = os.Stat(r.dataFilePath)
	if err != nil {
		r.items = make(map[string]*models.MenuItem)
		r.loaded = true
		return r.saveToFile()
	}

	file, err := os.Open(r.dataFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		r.items = make(map[string]*models.MenuItem)
		r.loaded = true
		return nil
	}

	items := []*models.MenuItem{}
	err = json.Unmarshal(data, &items)
	if err != nil {
		return err
	}

	r.items = make(map[string]*models.MenuItem)
	for _, item := range items {
		r.items[item.ID] = item
	}

	r.loaded = true
	r.logger.Debug("Loaded menu items from file", "count", len(r.items))
	return nil
}

func (r *MenuRepository) saveToFile() error {
	items := make([]*models.MenuItem, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}

	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(r.dataFilePath), 0755)
	if err != nil {
		return err
	}

	tempFile := r.dataFilePath + ".tmp"
	err = os.WriteFile(tempFile, data, 0644)
	if err != nil {
		return err
	}

	err = os.Rename(tempFile, r.dataFilePath)
	if err != nil {
		return err
	}

	r.logger.Debug("Save menu items to file", "count", len(items))
	return nil
}

func (r *MenuRepository) validateMenuItem(item *models.MenuItem) error {
	if item == nil {
		return errors.New("menu item cannot be nil")
	}
	if item.ID == "" {
		return errors.New("item ID cannot be empty")
	}
	if item.Name == "" {
		return errors.New("item name cannot be empty")
	}
	if item.Price < 0 {
		return errors.New("price cannot be negative")
	}

	if len(item.Ingredients) == 0 {
		return errors.New("menu item must have at least 1 ingridient")
	}
	for i, ingridient := range item.Ingredients {
		if ingridient.IngredientID == "" {
			return fmt.Errorf("ingridient %d: ID cannot be empty", i+1)
		}
		if ingridient.Quantity < 0 {
			return fmt.Errorf("ingridient %d: quantity must be positive", i+1)
		}
	}

	return nil
}

func (r *MenuRepository) backupFile() error {
	_, err := os.Stat(r.dataFilePath)
	if os.IsNotExist(err) {
		return nil
	}

	backupPath := r.dataFilePath + ".backup." + time.Now().Format("20060102_150405")

	data, err := os.ReadFile(r.dataFilePath)
	if err != nil {
		return err
	}
	err = os.WriteFile(backupPath, data, 0644)
	if err != nil {
		return err
	}

	r.logger.Debug("Created backup file", "backup_path", backupPath)
	return nil
}
