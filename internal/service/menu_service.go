package service

import (
	"fmt"
	"hot-coffee/internal/repositories"
	"hot-coffee/models"
	"hot-coffee/pkg/logger"
)

type CreateMenuItemRequest struct {
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Category    models.MenuCategory         `json:"category"`
	Price       float64                     `json:"price"`
	Available   bool                        `json:"available"`
	Ingridients []models.MenuItemIngredient `json:"ingridients"`
}

type UpdateMenuItemRequest struct {
	Name        *string                      `json:"name"`
	Description *string                      `json:"description"`
	Category    *models.MenuCategory         `json:"category"`
	Price       *float64                     `json:"price"`
	Available   *bool                        `json:"available"`
	Ingridients *[]models.MenuItemIngredient `json:"ingridients"`
}

type MenuServiceInterface interface {
	GetAllMenuItems() ([]*models.MenuItem, error)
	GetMenuItem(id string) (*models.MenuItem, error)
	CreateMenuItem(id string, req CreateMenuItemRequest) (*models.MenuItem, error)
	UpdateMenuItem(id string, req UpdateMenuItemRequest) error
	DeleteMenuItem(id string) error
	GetPopularItems() ([]*models.PopularItemAggregation, error)
}

type MenuService struct {
	menuRepo  repositories.MenuRepositoryInterface
	logger    *logger.Logger
}

func NewMenuService(menuRepo repositories.MenuRepositoryInterface, logger *logger.Logger) *MenuService {
	return &MenuService{
		menuRepo:  menuRepo,
		logger:    logger.WithComponent("menu_service"),
	}
}

// GetAllMenuItems retrieves all menu items
func (s *MenuService) GetAllMenuItems() ([]*models.MenuItem, error) {
	s.logger.Info("Fetching all menu items from repository")

	items, err := s.menuRepo.GetAll()
	if err != nil {
		s.logger.Error("Failed to get menu items from repository", "error", err)
		return nil, err
	}

	s.logger.Info("Fetched menu items", "count", len(items))
	return items, nil
}

// CreateMenuItem creates new menu item
func (s *MenuService) CreateMenuItem(id string, req CreateMenuItemRequest) (*models.MenuItem, error) {
	s.logger.Info("Creating menu item", "id", id, "name", req.Name, "price", req.Price)
	
	if err := s.validateCreateMenuItemData(req); err != nil {
		s.logger.Warn("Create failed: invalid data", "id", id, "error", err)
		return nil, err
	}

	item := &models.MenuItem{
		ID: id,
		Name: req.Name,
		Description: req.Description,
		Category: req.Category,
		Price: req.Price,
		Available: req.Available,
		Ingredients: req.Ingridients,
	}

	if err := s.menuRepo.Create(item); err != nil {
		s.logger.Error("Failed to create menu item in repository", "id", id, "error", err)
		return nil, err
	}

	s.logger.Info("Menu item created successfully", "id", id, "name", req.Name)
	return item, nil
}

// UpdateMenuItem updates existing menu item
func (s *MenuService) UpdateMenuItem(id string, req UpdateMenuItemRequest) error {
	s.logger.Info("Updating menu item", "id", id, "name", req.Name, "price", req.Price)
	
	if err := s.validateUpdateMenuItemData(req); err != nil {
		s.logger.Warn("Update failed: invalid data", "id", id, "error", err)
		return err
	}

	if _, err := s.menuRepo.GetByID(id); err != nil {
		s.logger.Error("Failed to get existing menu item", "id", id, "error", err)
		return err
	}

	item := &models.MenuItem{
		ID:id,
		Name: *req.Name,
		Description: *req.Description,
		Category: *req.Category,
		Price: *req.Price,
		Available: *req.Available,
		Ingredients: *req.Ingridients,
	}

	if err := s.menuRepo.Update(id, item); err != nil {
		s.logger.Error("Failed to update menu item", "id", id, "error", err)
		return err
	}

	s.logger.Info("Menu item updated successfully", "id", id, "name", req.Name)
	return nil
}

// DeleteMenuItem deletes menu item
func (s *MenuService) DeleteMenuItem(id string) error {
	s.logger.Info("Deleting menu item", "id", id)

	if _, err := s.menuRepo.GetByID(id); err != nil {
        s.logger.Warn("Menu item not found for deletion", "id", id, "error", err)
		return err
	}

	if err := s.menuRepo.Delete(id); err != nil {
        s.logger.Error("Failed to delete menu item from repository", "id", id, "error", err)
		return err
	}

    s.logger.Info("Menu item deleted successfully", "id", id)
	return nil
}

// GetMenuItem retrieves menu item by ID
func (s *MenuService) GetMenuItem(id string) (*models.MenuItem, error) {
	item, err := s.menuRepo.GetByID(id)
	if err != nil {
        s.logger.Warn("Menu item not found", "id", id, "error", err)
		return nil, err
	}

    s.logger.Info("Fetched menu item successfully", "id", id, "name", item.Name)
	return item, nil
}

// TODO: Implement GetPopularItems method - Get popular menu items
// - Call repository for popular items aggregation
// - Apply business logic for ranking
// - Log aggregation calculation
// func (s *MenuService) GetPopularItems() ([]*models.PopularItemAggregation, error)

func (s *MenuService) validateCreateMenuItemData(req CreateMenuItemRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Price < 0 {
		return fmt.Errorf("price must be non-negative")
	}
	if len(req.Ingridients) == 0 {
		return fmt.Errorf("menu item must have at least 1 ingridient")
	}

	for i, ingridient := range req.Ingridients {
		if ingridient.IngredientID == "" {
			return fmt.Errorf("ingridient %d: ID is required", i+1)
		}
		if ingridient.Quantity <= 0 {
			return fmt.Errorf("ingridient %d: quantity must be positive", i+1)
		}
	}

	return nil
}

func (s *MenuService) validateUpdateMenuItemData(req UpdateMenuItemRequest) error {
	if *req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if *req.Price < 0 {
		return fmt.Errorf("price must be non-negative")
	}
	if err := s.validateMenuCategory(*req.Category); err != nil {
		return err
	}
	if len(*req.Ingridients) == 0 {
		return fmt.Errorf("menu item must have at least 1 ingridient")
	}

	for i, ingridient := range *req.Ingridients {
		if ingridient.IngredientID == "" {
			return fmt.Errorf("ingridient %d: ID is required", i+1)
		}
		if ingridient.Quantity <= 0 {
			return fmt.Errorf("ingridient %d: quantity must be positive", i+1)
		}
	}

	return nil
}

func (s *MenuService) validateMenuCategory(category models.MenuCategory) error {
	switch category {
	case models.CategoryCoffee, models.CategoryDrink, models.CategoryPastry, models.CategorySandwich, models.CategoryTea:
		return nil
	default:
		return fmt.Errorf("invalid menu category: %s", category)
	}
}
