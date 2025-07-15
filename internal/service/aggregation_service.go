package service

import (
	"hot-coffee/internal/repositories"
	"hot-coffee/pkg/logger"
)

type AggregationServiceInterface interface {
	GetTotalSales() (*repositories.TotalSales, error)
	GetPopularItems() ([]repositories.PopularItem, error)
}

type AggregationService struct {
	aggregationRepo repositories.AggregationRepositoryInterface
	orderRepo       repositories.OrderRepositoryInterface
	menuRepo        repositories.MenuRepositoryInterface
	logger          *logger.Logger
}

func NewAggregationService(aggregationRepo repositories.AggregationRepositoryInterface, orderRepo repositories.OrderRepositoryInterface, menuRepo repositories.MenuRepositoryInterface, log *logger.Logger) *AggregationService {
	return &AggregationService{
		aggregationRepo: aggregationRepo,
		orderRepo:       orderRepo,
		menuRepo:        menuRepo,
		logger:          log.WithComponent("aggregation_repository"),
	}
}

func (s *AggregationService) GetTotalSales() (*repositories.TotalSales, error) {
	s.logger.Info("Getting total sales report")
	
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		s.logger.Error("Failed to get orders for sales report", "error", err)
		return nil, err
	}

	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		s.logger.Error("Failed to get menu items for sales report", "error", err)
		return nil, err
	}

	return s.aggregationRepo.CalculateTotalSales(orders, menuItems)
}

func (s *AggregationService) GetPopularItems() ([]repositories.PopularItem, error) {
	s.logger.Info("Getting popular itsme report")
	
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		s.logger.Error("Failed to get orders for popular items report", "error", err)
		return nil, err
	}

	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		s.logger.Error("Failed to get menu items for popular items report", "error", err)
		return nil, err
	}

	return s.aggregationRepo.CalculatePopularItems(orders, menuItems)
}
