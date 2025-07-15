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
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		// logger err
		return nil, err
	}

	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		// logger err
		return nil, err
	}

	return s.aggregationRepo.CalculateTotalSales(orders, menuItems)
}

func (s *AggregationService) GetPopularItems() ([]repositories.PopularItem, error) {
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		// logger err
		return nil, err
	}

	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		// logger err
		return nil, err
	}

	return s.aggregationRepo.CalculatePopularItems(orders, menuItems)
}
