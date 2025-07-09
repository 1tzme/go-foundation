# Hot Coffee - Three-Layered Architecture Documentation

## Project Overview

Hot Coffee is a coffee shop management system built with Go using a three-layered architecture pattern. The application provides HTTP endpoints for managing orders, menu items, and inventory with data persistence in JSON files.

## Architecture Overview

The project follows a three-layered architecture pattern:

1. **Presentation Layer (Handlers)** - HTTP request/response handling
2. **Business Logic Layer (Services)** - Core business logic and rules
3. **Data Access Layer (Repositories)** - Data storage and retrieval

## Project Structure

```
hot-coffee/
├── cmd/                            # Entry point of the application
│   └── main.go                     # Application bootstrap and HTTP server setup
├── internal/                       # Application code divided into layers
│   ├── handler/                    # Presentation layer handling HTTP requests
│   │   ├── order_handler.go        # Order-related HTTP endpoints
│   │   ├── menu_handler.go         # Menu-related HTTP endpoints
│   │   └── inventory_handler.go    # Inventory-related HTTP endpoints
│   ├── service/                    # Business logic layer with interfaces and implementations
│   │   ├── order_service.go        # Order business logic
│   │   ├── menu_service.go         # Menu business logic
│   │   └── inventory_service.go    # Inventory business logic
│   └── dal/                        # Data access layer with interfaces and implementations
│       ├── order_repository.go     # Order data operations
│       ├── menu_repository.go      # Menu data operations
│       └── inventory_repository.go # Inventory data operations
├── models/                         # Data models shared across layers
│   ├── order.go                    # Order-related data structures
│   ├── menu.go                     # Menu-related data structures
│   └── inventory.go                # Inventory-related data structures
├── data/                           # JSON data storage directory
│   ├── orders.json                 # Orders data persistence
│   ├── menu_items.json             # Menu items data persistence
│   └── inventory.json              # Inventory data persistence
├── go.mod                          # Go module definition
├── go.sum                          # Go module dependencies checksum
└── README.md                       # Project documentation
```

## Layer Responsibilities

### 1. Presentation Layer (Handlers)

**Location**: `internal/handler/`

**Responsibilities**:
- Handle HTTP requests and responses
- Parse input data and format output data
- Invoke appropriate methods from the Business Logic Layer
- Validate input data and return meaningful error messages
- Return appropriate HTTP status codes

**Implementation Details**:
- Organized by entity (orders, menu, inventory)
- Uses Gin framework for HTTP routing
- Handles JSON serialization/deserialization
- Implements proper error handling and logging

**Key Files**:
- `order_handler.go` - Order CRUD operations
- `menu_handler.go` - Menu item management
- `inventory_handler.go` - Inventory management

### 2. Business Logic Layer (Services)

**Location**: `internal/service/`

**Responsibilities**:
- Implement core business logic and rules
- Define interfaces for services to promote decoupling
- Perform data processing and validation
- Handle aggregations and computations
- Coordinate between handlers and repositories

**Implementation Details**:
- Service interfaces for testability and flexibility
- Business rule validation
- Data aggregation methods
- Error handling and logging
- Independent and testable components

**Key Files**:
- `order_service.go` - Order business logic and validation
- `menu_service.go` - Menu management and availability logic
- `inventory_service.go` - Inventory tracking and low-stock alerts

### 3. Data Access Layer (Repositories)

**Location**: `internal/dal/`

**Responsibilities**:
- Manage data storage and retrieval operations
- Interact with JSON files to persist and read data
- Ensure data integrity and consistency
- Provide interfaces for flexibility

**Implementation Details**:
- Repository interfaces for each entity
- JSON file-based persistence
- Thread-safe operations using mutexes
- Data validation and error handling
- Atomic file operations

**Key Files**:
- `order_repository.go` - Order data persistence
- `menu_repository.go` - Menu item data persistence
- `inventory_repository.go` - Inventory data persistence

## API Endpoints

### Order Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/orders` | Create a new order |
| GET | `/api/v1/orders` | Get all orders |
| GET | `/api/v1/orders/:id` | Get order by ID |
| PUT | `/api/v1/orders/:id/status` | Update order status |
| GET | `/api/v1/orders/aggregations/sales` | Get total sales |

### Menu Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/menu` | Get all menu items |
| POST | `/api/v1/menu` | Create new menu item |
| PUT | `/api/v1/menu/:id` | Update menu item |
| DELETE | `/api/v1/menu/:id` | Delete menu item |
| GET | `/api/v1/menu/aggregations/popular` | Get popular menu items |

### Inventory Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/inventory` | Get all inventory items |
| PUT | `/api/v1/inventory/:id` | Update inventory item |
| GET | `/api/v1/inventory/low-stock` | Get low stock items |

## Data Models

### Order Model

```go
type Order struct {
    ID          string      `json:"id"`
    CustomerID  string      `json:"customer_id"`
    Items       []OrderItem `json:"items"`
    TotalAmount float64     `json:"total_amount"`
    Status      OrderStatus `json:"status"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
}

type OrderStatus string
const (
    OrderStatusPending    OrderStatus = "pending"
    OrderStatusPreparing  OrderStatus = "preparing"
    OrderStatusReady      OrderStatus = "ready"
    OrderStatusCompleted  OrderStatus = "completed"
    OrderStatusCancelled  OrderStatus = "cancelled"
)
```

### Menu Model

```go
type MenuItem struct {
    ID          string       `json:"id"`
    Name        string       `json:"name"`
    Description string       `json:"description"`
    Category    MenuCategory `json:"category"`
    Price       float64      `json:"price"`
    Available   bool         `json:"available"`
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
}

type MenuCategory string
const (
    CategoryCoffee    MenuCategory = "coffee"
    CategoryTea       MenuCategory = "tea"
    CategoryPastry    MenuCategory = "pastry"
    CategorySandwich  MenuCategory = "sandwich"
    CategoryDrink     MenuCategory = "drink"
)
```

### Inventory Model

```go
type InventoryItem struct {
    ID           string    `json:"id"`
    Name         string    `json:"name"`
    Description  string    `json:"description"`
    Quantity     int       `json:"quantity"`
    MinThreshold int       `json:"min_threshold"`
    Unit         string    `json:"unit"`
    LastUpdated  time.Time `json:"last_updated"`
}
```

## Data Persistence

### JSON File Storage

Data is persisted in JSON files located in the `data/` directory:

- `data/orders.json` - Contains all order records
- `data/menu_items.json` - Contains all menu item records
- `data/inventory.json` - Contains all inventory item records

### File Operations

- **Atomic Writes**: Files are written atomically to prevent corruption
- **Backup**: Previous versions are backed up before updates
- **Thread Safety**: Concurrent access is handled with mutexes
- **Error Recovery**: Graceful handling of file system errors

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DATA_DIR` | `./data` | Data storage directory |
| `LOG_LEVEL` | `info` | Logging level |

### Server Configuration

The application starts an HTTP server with configurable settings:

```go
// Default configuration
config := Config{
    Port:    getEnv("PORT", "8080"),
    DataDir: getEnv("DATA_DIR", "./data"),
    LogLevel: getEnv("LOG_LEVEL", "info"),
}
```

## Logging

The application uses Go's `log/slog` package for structured logging:

### Log Levels

- **DEBUG**: Detailed debugging information
- **INFO**: General information about application flow
- **WARN**: Warning conditions that don't prevent operation
- **ERROR**: Error conditions that require attention

### Log Events

- HTTP request/response logging
- Business logic operations
- Data access operations
- Error conditions
- Performance metrics

### Example Log Entries

```json
{
  "time": "2025-07-10T10:30:00Z",
  "level": "INFO",
  "msg": "Order created successfully",
  "order_id": "ord-123",
  "customer_id": "cust-456",
  "total_amount": 15.50
}

{
  "time": "2025-07-10T10:31:00Z",
  "level": "ERROR",
  "msg": "Failed to update inventory",
  "item_id": "inv-789",
  "error": "file not found"
}
```

## Error Handling

### HTTP Status Codes

- `200 OK` - Successful operation
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid input data
- `404 Not Found` - Resource not found
- `409 Conflict` - Business rule violation
- `500 Internal Server Error` - Server error

### Error Response Format

```json
{
  "error": {
    "code": "INVALID_INPUT",
    "message": "Order must contain at least one item",
    "details": {
      "field": "items",
      "value": "[]"
    }
  }
}
```

## Aggregations

### Sales Aggregations

- **Total Sales**: Sum of all completed orders
- **Sales by Period**: Daily, weekly, monthly sales
- **Sales by Category**: Revenue breakdown by menu category

### Menu Aggregations

- **Popular Items**: Most frequently ordered items
- **Revenue by Item**: Total revenue per menu item
- **Category Performance**: Sales performance by category

### Inventory Aggregations

- **Low Stock Items**: Items below minimum threshold
- **Inventory Value**: Total value of current inventory
- **Usage Patterns**: Consumption rates by item

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

Clone the repository:
```bash
git clone <repository-url>
cd hot-coffee
```

### Usage

```bash
$ ./hot-coffee --help
Coffee Shop Management System

Usage:
  hot-coffee [--port <N>] [--dir <S>] 
  hot-coffee --help

Options:
  --help       Show this screen.
  --port N     Port number.
  --dir S      Path to the data directory.
```

### Development

1. Follow the layered architecture pattern
2. Add new features by creating corresponding files in each layer
3. Ensure proper error handling and logging
4. Write unit tests for business logic
5. Document API changes

## Best Practices

### Code Organization

- Keep layers separate and well-defined
- Use interfaces for dependency injection
- Follow Go naming conventions
- Implement proper error handling

### Data Management

- Validate data at service layer
- Use atomic file operations
- Implement proper backup strategies
- Handle concurrent access safely

### API Design

- Use RESTful conventions
- Return consistent error formats
- Implement proper HTTP status codes
- Version your APIs

### Testing

- Write unit tests for business logic
- Mock dependencies for isolation
- Test error conditions
- Implement integration tests

```
This documentation provides a comprehensive guide for understanding, developing, and maintaining the Hot Coffee application using the three-layered architecture pattern.
```
