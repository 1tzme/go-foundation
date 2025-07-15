package repositories

import (
	"hot-coffee/models"
	"hot-coffee/pkg/logger"
	"sort"
)

type AggregationRepositoryInterface interface{
	CalculateTotalSales(orders []*models.Order, menuItems []*models.MenuItem) (*TotalSales, error)
	CalculatePopularItems(orders []*models.Order, menuItems []*models.MenuItem) ([]PopularItem, error)
}

type TotalSales struct {
	TotalRevenue float64    `json:"total_revenue"`
	ItemSales    []ItemSale `json:"item_sales"`
}

type ItemSale struct {
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	QuantitySold int     `json:"quantity_sold"`
	TotalValue   float64 `json:"total_value"`
}

type PopularItem struct {
	models.MenuItem
	SalesCount int `json:"sales_count"`
}

type AggregationRepository struct {
	logger *logger.Logger
}

func NewAggregationRepository(log *logger.Logger) *AggregationRepository {
	return &AggregationRepository{
		logger: log.WithComponent("aggregation_repository"),
	}
}

// CalculateTotalSales - retrieves calculated total sales report
func (r *AggregationRepository) CalculateTotalSales(orders []*models.Order, menuItems []*models.MenuItem) (*TotalSales, error) {
	r.logger.Info("Calculating total sales from provided data")

	
	menuMap := make(map[string]*models.MenuItem)
	for _, item := range menuItems {
		menuMap[item.ID] = item
	}

	report := &TotalSales{
		ItemSales: make([]ItemSale, 0),
	}
	itemSalesMap := make(map[string]*ItemSale)

	for _, order := range orders {
		if order.Status != "closed" {
			continue
		}

		for _, orderItem := range order.Items {
			menuItem, ok := menuMap[orderItem.ProductID]
			if !ok {
				r.logger.Warn("Product ID from order not found in menu", "product_id", orderItem.ProductID, "order_id", order.ID)
				continue
			}

			itemValue := menuItem.Price * float64(orderItem.Quantity)
			report.TotalRevenue += itemValue

			sale, exists := itemSalesMap[orderItem.ProductID]
			if exists {
				sale.QuantitySold += orderItem.Quantity
				sale.TotalValue += itemValue
			} else {
				itemSalesMap[orderItem.ProductID] = &ItemSale{
					ProductID: orderItem.ProductID,
					ProductName: menuItem.Name,
					QuantitySold: orderItem.Quantity,
					TotalValue: itemValue,
				}
			}
		}
	}

	for _, sale := range itemSalesMap {
		report.ItemSales = append(report.ItemSales, *sale)
	}
	
	sort.Slice(report.ItemSales, func(i, j int) bool {
		return report.ItemSales[i].ProductName < report.ItemSales[j].ProductName
	})
	return report, nil
}

func (r *AggregationRepository) CalculatePopularItems(orders []*models.Order, menuItems []*models.MenuItem) ([]PopularItem, error) {
	r.logger.Info("Calculating popular itsme from closed orders")
	
	salesCount := make(map[string]int)
	for _, order := range orders {
		if order.Status != "closed" {
			continue
		}
		for _, item := range order.Items {
			salesCount[item.ProductID] += item.Quantity
		}
	}

	popularItems := []PopularItem{}
	for _, menuItem := range menuItems {
		popularItems = append(popularItems, PopularItem{
			MenuItem: *menuItem,
			SalesCount: salesCount[menuItem.ID],
		})
	}

	sort.Slice(popularItems, func(i, j int) bool {
		return popularItems[i].SalesCount > popularItems[j].SalesCount
	})

	return popularItems, nil
}
