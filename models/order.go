package models

// TODO: Add import when implementing time-based fields:
// import "time"

type Order struct {
    ID           string       `json:"order_id"`
    CustomerName string       `json:"customer_name"`
    Items        []OrderItem  `json:"items"`
    Status       string       `json:"status"`
    CreatedAt    string       `json:"created_at"`
    // TODO: Add additional fields based on README spec:
    // CustomerID  string      `json:"customer_id"`
    // TotalAmount float64     `json:"total_amount"`
    // UpdatedAt   time.Time   `json:"updated_at"`
}

type OrderItem struct {
    ProductID string `json:"product_id"`
    Quantity  int    `json:"quantity"`
}

// TODO: Add OrderStatus enum based on README spec
// type OrderStatus string
// const (
//     OrderStatusPending    OrderStatus = "pending"
//     OrderStatusPreparing  OrderStatus = "preparing"
//     OrderStatusReady      OrderStatus = "ready"
//     OrderStatusCompleted  OrderStatus = "completed"
//     OrderStatusCancelled  OrderStatus = "cancelled"
// )

// TODO: Add aggregation models based on README spec
// type SalesAggregation struct {
//     TotalSales      float64                    `json:"total_sales"`
//     SalesByPeriod   map[string]float64        `json:"sales_by_period"`
//     SalesByCategory map[string]float64        `json:"sales_by_category"`
//     OrderCount      int                       `json:"order_count"`
//     AverageOrderValue float64                 `json:"average_order_value"`
//     LastUpdated     time.Time                 `json:"last_updated"`
// }