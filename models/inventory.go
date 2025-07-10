package models

import "time"

type InventoryItem struct {
	IngredientID string    `json:"ingredient_id"`
	Name         string    `json:"name"`
	Quantity     float64   `json:"quantity"`
	Unit         string    `json:"unit"`
	ID           string    `json:"id"`
	Description  string    `json:"description"`
	MinThreshold float64   `json:"min_threshold"`
	LastUpdated  time.Time `json:"last_updated"`
}

// TODO: Add aggregation models based on README spec
type InventoryValueAggregation struct {
	TotalValue      float64            `json:"total_value"`
	ValueByCategory map[string]float64 `json:"value_by_category"`
	ItemCount       int                `json:"item_count"`
	LowStockCount   int                `json:"low_stock_count"`
	LastCalculated  time.Time          `json:"last_calculated"`
}

// TODO: Add low stock alert model
type LowStockAlert struct {
	ItemID          string    `json:"item_id"`
	ItemName        string    `json:"item_name"`
	CurrentQuantity int       `json:"current_quantity"`
	MinThreshold    int       `json:"min_threshold"`
	AlertLevel      string    `json:"alert_level"` // "low", "critical"
	LastUpdated     time.Time `json:"last_updated"`
}
