package service

// TODO: Add imports when implementing:
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

// TODO: Implement MenuService interface
type MenuServiceInterface interface {
	GetAllMenuItems() ([]*models.MenuItem, error)
	GetMenuItem(id string) ()
	CreateMenuItem(id string, req CreateMenuItemRequest) (*models.MenuItem, error)
	UpdateMenuItem(id string, req UpdateMenuItemRequest) error
	DeleteMenuItem(id string) error
	GetPopularItems() ([]*models.PopularItemAggregation, error)
}

// TODO: Implement MenuService struct
type MenuService struct {
	menuRepo  repositories.MenuRepositoryInterface
	logger    *logger.Logger
}

// TODO: Implement constructor with logger injection
func NewMenuService(menuRepo repositories.MenuRepositoryInterface, logger *logger.Logger) *MenuService {
	return &MenuService{
		menuRepo:  menuRepo,
		logger:    logger.WithComponent("menu_service"),
	}
}

// TODO: Implement GetAllMenuItems method - Retrieve all menu items
// - Call repository to get all menu items
// - Apply business logic for availability
// - Log retrieval event
// func (s *MenuService) GetAllMenuItems() ([]*models.MenuItem, error)

// TODO: Implement CreateMenuItem method - Create new menu item
// - Validate menu item data (name, price > 0, valid category)
// - Check for duplicate names
// - Set created timestamp
// - Call repository to create
// - Log business event
// func (s *MenuService) CreateMenuItem(id string, req CreateMenuItemRequest) (*models.MenuItem, error)

// TODO: Implement UpdateMenuItem method - Update existing menu item
// - Validate menu item exists
// - Validate updated data
// - Set updated timestamp
// - Call repository to update
// - Log business event
// func (s *MenuService) UpdateMenuItem(id string, req MenuItemRequest) error

// TODO: Implement DeleteMenuItem method - Delete menu item
// - Validate menu item exists
// - Check if item is referenced in active orders
// - Call repository to delete
// - Log business event
// func (s *MenuService) DeleteMenuItem(id string) error

// TODO: Implement GetPopularItems method - Get popular menu items
// - Call repository for popular items aggregation
// - Apply business logic for ranking
// - Log aggregation calculation
// func (s *MenuService) GetPopularItems() ([]*models.PopularItemAggregation, error)

// TODO: Implement private business logic methods
// - validateMenuItemData(req MenuItemRequest) error - Validate menu item
// - checkDuplicateName(name string, excludeID string) error - Check for duplicate names
// - isMenuItemInActiveOrders(id string) (bool, error) - Check if item is in active orders
// - validateMenuCategory(category models.MenuCategory) error - Validate category
// - validatePrice(price float64) error - Validate price is positive

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
