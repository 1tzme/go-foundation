# Logger Placement Guide for Hot Coffee Project

## Where Should Loggers Be Placed?

### 1. **Centralized Logger Package** (`pkg/logger/`)
- **Purpose**: Provide consistent logging configuration across the application
- **Location**: `pkg/logger/logger.go`
- **Responsibility**: Configure slog with proper levels, formats, and handlers

### 2. **Application Bootstrap** (`cmd/main.go`)
- **Purpose**: Initialize logger and inject it into all components
- **Pattern**: Dependency injection from the top level
- **Example**:
```go
func main() {
    // Initialize logger
    appLogger := logger.New(logger.Config{
        Level: "info",
        Format: "json",
    })
    
    // Inject into repositories
    orderRepo := dal.NewOrderRepository(appLogger)
    
    // Inject into services
    orderService := service.NewOrderService(orderRepo, appLogger)
    
    // Inject into handlers
    orderHandler := handler.NewOrderHandler(orderService, appLogger)
}
```

### 3. **Data Access Layer** (`internal/dal/`)
- **Constructor Pattern**: Accept logger in constructor
- **Context**: Add component context (`"component", "order_repository"`)
- **Log Events**:
  - File I/O operations
  - Data persistence operations
  - Error conditions
  - Performance metrics

### 4. **Business Logic Layer** (`internal/service/`)
- **Constructor Pattern**: Accept logger in constructor
- **Context**: Add component context (`"component", "order_service"`)
- **Log Events**:
  - Business operations
  - Validation failures
  - State changes
  - Business rule violations

### 5. **Presentation Layer** (`internal/handler/`)
- **Constructor Pattern**: Accept logger in constructor
- **Context**: Add component context (`"component", "order_handler"`)
- **Log Events**:
  - HTTP requests/responses
  - Input validation errors
  - Status codes returned
  - Request processing time

## Best Practices Summary

### ✅ DO:
- **Use dependency injection** - Pass logger through constructors
- **Add context** - Use `logger.WithContext()` for component identification
- **Log structured data** - Use key-value pairs
- **Log at appropriate levels** - DEBUG, INFO, WARN, ERROR
- **Include relevant IDs** - order_id, customer_id, etc.

### ❌ DON'T:
- **Use global loggers** - Avoid `slog.Default()` in business logic
- **Log sensitive data** - Passwords, tokens, personal info
- **Over-log** - Avoid excessive DEBUG logs in production
- **Log without context** - Always include relevant metadata

## Logger Flow Example

```
HTTP Request → Handler (logs request) 
    ↓
Service (logs business operation)
    ↓  
Repository (logs data operation)
    ↓
JSON File (actual persistence)
```

Each layer logs its specific concerns while maintaining the same logger instance with added context for traceability.

## Implementation Pattern

Every component should follow this pattern:

```go
type ComponentName struct {
    dependencies *SomeDependency
    logger       *logger.Logger  // Always inject
}

func NewComponentName(deps *SomeDependency, logger *logger.Logger) *ComponentName {
    return &ComponentName{
        dependencies: deps,
        logger:       logger.WithContext("component", "component_name"),
    }
}

func (c *ComponentName) SomeMethod() {
    c.logger.Info("Starting operation", "param", value)
    // ... operation logic
    c.logger.Info("Operation completed", "result", result)
}
```

This ensures consistent, traceable, and maintainable logging throughout the application.
