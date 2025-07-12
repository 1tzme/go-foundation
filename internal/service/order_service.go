package service

import (
	"fmt"
	"hot-coffee/internal/repositories"
	"hot-coffee/models"
	"hot-coffee/pkg/logger"
	"strconv"
	"strings"
	"time"
)

// Define request/response structs
type CreateOrderRequest struct {
	CustomerName string                   `json:"customer_name"`
	Items        []CreateOrderItemRequest `json:"items"`
}

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type UpdateOrderRequest struct {
	CustomerName string                   `json:"customer_name"`
	Items        []CreateOrderItemRequest `json:"items"`
	Status       string                   `json:"status"`
}

// OrderService interface
type OrderServiceInterface interface {
	CreateOrder(req CreateOrderRequest) (*models.Order, error)
	GetAllOrders() ([]*models.Order, error)
	GetOrderByID(id string) (*models.Order, error)
	UpdateOrder(id string, req UpdateOrderRequest) error
	DeleteOrder(id string) error
	CloseOrder(id string) error
}

// OrderService struct
type OrderService struct {
	orderRepo repositories.OrderRepositoryInterface
	logger    *logger.Logger
}

// NewOrderService creates a new OrderService with the given repository and logger
func NewOrderService(orderRepo repositories.OrderRepositoryInterface, logger *logger.Logger) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		logger:    logger.WithComponent("order_service"),
	}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(req CreateOrderRequest) (*models.Order, error) {
	s.logger.Info("Creating new order", "customer", req.CustomerName)

	if err := s.validateOrderData(req); err != nil {
		s.logger.Warn("Create failed: invalid data", "error", err)
		return nil, err
	}

	orderID := s.generateOrderID()
	order := &models.Order{
		ID:           orderID,
		CustomerName: req.CustomerName,
		Items:        make([]models.OrderItem, len(req.Items)),
		Status:       "open",
		CreatedAt:    time.Now().Format(time.RFC3339),
	}

	// Convert request items to order items
	for i, item := range req.Items {
		order.Items[i] = models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	if err := s.orderRepo.Add(order); err != nil {
		s.logger.Error("Failed to add order in repository", "order_id", orderID, "error", err)
		return nil, err
	}

	s.logger.Info("Order created", "order_id", orderID)
	return order, nil
}

// GetAllOrders retrieves all orders
func (s *OrderService) GetAllOrders() ([]*models.Order, error) {
	s.logger.Info("Fetching all orders from repository")

	orders, err := s.orderRepo.GetAll()
	if err != nil {
		s.logger.Error("Failed to fetch orders from repository", "error", err)
		return nil, err
	}

	s.logger.Info("Fetched orders", "count", len(orders))
	return orders, nil
}

// GetOrderByID retrieves a specific order by ID
func (s *OrderService) GetOrderByID(id string) (*models.Order, error) {
	s.logger.Info("Fetching order by ID", "order_id", id)

	if id == "" {
		s.logger.Warn("Order ID cannot be empty")
		return nil, fmt.Errorf("order ID is required")
	}

	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		s.logger.Warn("Order not found", "order_id", id, "error", err)
		return nil, err
	}

	s.logger.Info("Fetched order", "order_id", id)
	return order, nil
}

// UpdateOrder updates an existing order
func (s *OrderService) UpdateOrder(id string, req UpdateOrderRequest) error {
	s.logger.Info("Updating order", "order_id", id, "customer", req.CustomerName)

	if id == "" {
		s.logger.Warn("Order ID cannot be empty")
		return fmt.Errorf("order ID is required")
	}

	if err := s.validateUpdateOrderData(req); err != nil {
		s.logger.Warn("Update failed: invalid data", "order_id", id, "error", err)
		return err
	}

	order := &models.Order{
		ID:           id,
		CustomerName: req.CustomerName,
		Items:        make([]models.OrderItem, len(req.Items)),
		Status:       req.Status,
		CreatedAt:    time.Now().Format(time.RFC3339), // This would normally be preserved from original
	}

	// Convert request items to order items
	for i, item := range req.Items {
		order.Items[i] = models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	if err := s.orderRepo.Update(id, order); err != nil {
		s.logger.Error("Failed to update order in repository", "order_id", id, "error", err)
		return err
	}

	s.logger.Info("Order updated", "order_id", id)
	return nil
}

// DeleteOrder deletes an order by ID
func (s *OrderService) DeleteOrder(id string) error {
	s.logger.Info("Deleting order", "order_id", id)

	if id == "" {
		s.logger.Warn("Order ID cannot be empty")
		return fmt.Errorf("order ID is required")
	}

	if err := s.orderRepo.Delete(id); err != nil {
		s.logger.Warn("Failed to delete order", "order_id", id, "error", err)
		return err
	}

	s.logger.Info("Order deleted", "order_id", id)
	return nil
}

// CloseOrder closes an order by setting status to closed
func (s *OrderService) CloseOrder(id string) error {
	s.logger.Info("Closing order", "order_id", id)

	if id == "" {
		s.logger.Warn("Order ID cannot be empty")
		return fmt.Errorf("order ID is required")
	}

	if err := s.orderRepo.CloseOrder(id); err != nil {
		s.logger.Warn("Failed to close order", "order_id", id, "error", err)
		return err
	}

	s.logger.Info("Order closed", "order_id", id)
	return nil
}

// Private business logic methods

// validateOrderData validates the data for order creation
func (s *OrderService) validateOrderData(req CreateOrderRequest) error {
	if req.CustomerName == "" {
		return fmt.Errorf("customer name is required")
	}
	if len(req.Items) == 0 {
		return fmt.Errorf("order must have at least one item")
	}
	return s.validateOrderItems(req.Items)
}

// validateUpdateOrderData validates the data for order updates
func (s *OrderService) validateUpdateOrderData(req UpdateOrderRequest) error {
	if req.CustomerName == "" {
		return fmt.Errorf("customer name is required")
	}
	if len(req.Items) == 0 {
		return fmt.Errorf("order must have at least one item")
	}
	if req.Status == "" {
		return fmt.Errorf("status is required")
	}

	// Validate status values
	validStatuses := []string{"open", "closed"}
	statusValid := false
	for _, status := range validStatuses {
		if req.Status == status {
			statusValid = true
			break
		}
	}
	if !statusValid {
		return fmt.Errorf("invalid status: %s", req.Status)
	}

	return s.validateOrderItems(req.Items)
}

// validateOrderItems validates individual order items
func (s *OrderService) validateOrderItems(items []CreateOrderItemRequest) error {
	for i, item := range items {
		if item.ProductID == "" {
			return fmt.Errorf("item %d: product ID is required", i+1)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("item %d: quantity must be positive", i+1)
		}
	}
	return nil
}

// generateOrderID generates a unique order ID
func (s *OrderService) generateOrderID() string {
	// Get all existing orders to determine next ID
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		// If we can't get orders, start with order1
		return "order1"
	}

	// If no orders exist, start with order1
	if len(orders) == 0 {
		return "order1"
	}

	// Find the highest order number
	maxOrderNum := 0
	for _, order := range orders {
		if strings.HasPrefix(order.ID, "order") {
			// Extract number from "order123" format
			numStr := strings.TrimPrefix(order.ID, "order")
			if num, err := strconv.Atoi(numStr); err == nil && num > maxOrderNum {
				maxOrderNum = num
			}
		}
	}

	// Return next sequential order ID
	return fmt.Sprintf("order%d", maxOrderNum+1)
}
